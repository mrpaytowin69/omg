package resiphost

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"opensvc.com/opensvc/core/actioncontext"
	"opensvc.com/opensvc/core/actionrollback"
	"opensvc.com/opensvc/core/provisioned"
	"opensvc.com/opensvc/core/resource"
	"opensvc.com/opensvc/core/status"
	"opensvc.com/opensvc/util/fqdn"
	"opensvc.com/opensvc/util/hostname"
	"opensvc.com/opensvc/util/netif"

	"github.com/go-ping/ping"
)

const (
	tagNonRouted = "nonrouted"
)

type (
	T struct {
		resource.T

		// config
		IpName        string   `json:"ipname"`
		IpDev         string   `json:"ipdev"`
		Netmask       string   `json:"netmask"`
		Network       string   `json:"network"`
		Gateway       string   `json:"gateway"`
		WaitDNS       bool     `json:"wait_dns"`
		DNSUpdate     bool     `json:"dns_update"`
		DNSNameSuffix string   `json:"dns_name_suffix"`
		Provisioner   string   `json:"provisioner"`
		CheckCarrier  bool     `json:"check_carrier"`
		Alias         bool     `json:"alias"`
		Expose        []string `json:"expose"`

		// cache
		_ipaddr net.IP
		_ipmask net.IPMask
		_ipnet  *net.IPNet
	}

	Addrs []net.Addr
)

func New() resource.Driver {
	t := &T{}
	return t
}

func (t *T) StatusInfo() map[string]interface{} {
	netmask, _ := t.ipmask().Size()
	data := make(map[string]interface{})
	data["ipaddr"] = t.ipaddr()
	data["ipdev"] = t.IpDev
	data["netmask"] = netmask
	return data
}

func (t *T) Start(ctx context.Context) error {
	if initialStatus := t.Status(ctx); initialStatus == status.Up {
		t.Log().Info().Msgf("%s is already up on %s", t.IpName, t.IpDev)
		return nil
	}
	if err := t.start(); err != nil {
		return err
	}
	actionrollback.Register(ctx, func() error {
		return t.stop()
	})
	if err := t.arpAnnounce(); err != nil {
		return err
	}
	return nil
}

func (t *T) Stop(ctx context.Context) error {
	if initialStatus := t.Status(ctx); initialStatus == status.Down {
		t.Log().Info().Msgf("%s is already down on %s", t.IpName, t.IpDev)
		return nil
	}
	if err := t.stop(); err != nil {
		return err
	}
	return nil
}

func (t *T) Status(ctx context.Context) status.T {
	var (
		i       *net.Interface
		err     error
		addrs   Addrs
		carrier bool
	)
	ip := t.ipaddr()
	if t.IpName == "" {
		t.StatusLog().Warn("ipname not set")
		return status.NotApplicable
	}
	if t.IpDev == "" {
		t.StatusLog().Warn("ipdev not set")
		return status.NotApplicable
	}
	if i, err = t.netInterface(); err != nil {
		if fmt.Sprint(err.(*net.OpError).Unwrap()) == "no such network interface" {
			t.StatusLog().Warn("interface %s not found", t.IpDev)
		} else {
			t.StatusLog().Error("%s", err)
		}
		return status.Down
	}
	if t.CheckCarrier {
		if carrier, err = t.hasCarrier(); err == nil && carrier == false {
			t.StatusLog().Error("interface %s no-carrier.", t.IpDev)
			return status.Down
		}
	}
	if addrs, err = i.Addrs(); err != nil {
		t.StatusLog().Error("%s", err)
		return status.Down
	}
	if !Addrs(addrs).Has(ip) {
		t.Log().Debug().Msg("ip not found on intf")
		return status.Down
	}
	return status.Up
}

func (t T) Label() string {
	return fmt.Sprintf("%s %s", t.ipnet(), t.IpDev)
}

func (t *T) Provision(ctx context.Context) error {
	return nil
}

func (t *T) Unprovision(ctx context.Context) error {
	return nil
}

func (t T) Provisioned() (provisioned.T, error) {
	return provisioned.NotApplicable, nil
}

func (t T) Abort(ctx context.Context) bool {
	if t.Tags.Has(tagNonRouted) || t.IsActionDisabled() {
		return false // let start fail with an explicit error message
	}
	if t.ipaddr() == nil {
		return false // let start fail with an explicit error message
	}
	if initialStatus := t.Status(ctx); initialStatus == status.Up {
		return false // let start fail with an explicit error message
	}
	if t.CheckCarrier {
		if carrier, err := t.hasCarrier(); err == nil && carrier == false && !actioncontext.IsForce(ctx) {
			t.Log().Error().Msgf("interface %s no-carrier.", t.IpDev)
			return true
		}
	}
	if t.abortPing() {
		return true
	}
	return false
}

func (t T) hasCarrier() (bool, error) {
	return netif.HasCarrier(t.IpDev)
}

func (t T) abortPing() bool {
	ip := t.ipaddr()
	pinger, err := ping.NewPinger(ip.String())
	if err != nil {
		t.Log().Error().Err(err).Msg("abort: ping")
		return true
	}
	pinger.Count = 5
	pinger.Timeout = 5 * time.Second
	pinger.Interval = time.Second
	t.Log().Info().Msgf("checking %s availability (5s)", ip)
	pinger.Run()
	return pinger.Statistics().PacketsRecv > 0
}

func (t T) ipnet() *net.IPNet {
	if t._ipnet != nil {
		return t._ipnet
	}
	t._ipnet = t.getIPNet()
	return t._ipnet
}

func (t T) ipaddr() net.IP {
	if t._ipaddr != nil {
		return t._ipaddr
	}
	t._ipaddr = t.getIPAddr()
	return t._ipaddr
}

func (t *T) ipmask() net.IPMask {
	if t._ipmask != nil {
		return t._ipmask
	}
	t._ipmask = t.getIPMask()
	return t._ipmask
}

func (t *T) getIPNet() *net.IPNet {
	return &net.IPNet{
		IP:   t.ipaddr(),
		Mask: t.ipmask(),
	}
}

func (t *T) getIPMask() net.IPMask {
	ip := t.ipaddr()
	bits := getIPBits(ip)
	if m, err := parseCIDRMask(t.Netmask, bits); err == nil {
		return m
	}
	if m, err := parseDottedMask(t.Netmask); err == nil {
		return m
	}
	// fallback to the mask of the first found ip on the intf
	if m, err := t.defaultMask(); err == nil {
		return m
	}
	return nil
}

func (t *T) defaultMask() (net.IPMask, error) {
	intf, err := t.netInterface()
	if err != nil {
		return nil, err
	}
	addrs, err := intf.Addrs()
	if err != nil {
		return nil, err
	}
	if len(addrs) == 0 {
		return nil, fmt.Errorf("no addr to guess mask from")
	}
	_, net, err := net.ParseCIDR(addrs[0].String())
	if err != nil {
		return nil, err
	}
	return net.Mask, nil
}

func (t T) getIPAddr() net.IP {
	switch {
	case fqdn.IsValid(t.IpName) || hostname.IsValid(t.IpName):
		var (
			l   []net.IP
			err error
		)
		l, err = net.LookupIP(t.IpName)
		if err != nil {
			t.Log().Error().Err(err)
			return nil
		}
		n := len(l)
		switch n {
		case 0:
			t.Log().Error().Msgf("ipname %s is unresolvable", t.IpName)
		case 1:
			// ok
		default:
			t.Log().Debug().Msgf("ipname %s is resolvables to %d address. Using the first.", t.IpName, n)
		}
		return l[0]
	default:
		return net.ParseIP(t.IpName)
	}
}

func (t T) netInterface() (*net.Interface, error) {
	return net.InterfaceByName(t.IpDev)
}

func (t Addrs) Has(ip net.IP) bool {
	for _, addr := range t {
		listIP, _, _ := net.ParseCIDR(addr.String())
		if ip.Equal(listIP) {
			return true
		}
	}
	return false
}

func parseCIDRMask(s string, bits int) (net.IPMask, error) {
	if bits == 0 {
		return nil, errors.New("invalid bits: 0")
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return nil, fmt.Errorf("invalid element in dotted mask: %s", err)
	}
	return net.CIDRMask(i, bits), nil
}

func parseDottedMask(s string) (net.IPMask, error) {
	m := []byte{}
	l := strings.Split(s, ".")
	if len(l) != 4 {
		return nil, errors.New("invalid number of elements in dotted mask")
	}
	for _, e := range l {
		i, err := strconv.Atoi(e)
		if err != nil {
			return nil, fmt.Errorf("invalid element in dotted mask: %s", err)
		}
		m = append(m, byte(i))
	}
	return m, nil
}

func ipv4MaskString(m []byte) string {
	if len(m) != 4 {
		panic("ipv4Mask: len must be 4 bytes")
	}

	return fmt.Sprintf("%d.%d.%d.%d", m[0], m[1], m[2], m[3])
}

func getIPBits(ip net.IP) (bits int) {
	switch {
	case ip.To4() != nil:
		bits = 32
	case ip.To16() != nil:
		bits = 128
	}
	return
}

func (t T) arpAnnounce() error {
	ip := t.ipaddr()
	if ip.IsLoopback() {
		t.Log().Debug().Msgf("skip arp announce on loopback address %s", ip)
		return nil
	}
	if ip.IsLinkLocalUnicast() {
		t.Log().Debug().Msgf("skip arp announce on link local unicast address %s", ip)
		return nil
	}
	if ip.To4() == nil {
		t.Log().Debug().Msgf("skip arp announce on non-ip4 address %s", ip)
		return nil
	}
	if i, err := t.netInterface(); err == nil && i.Flags&net.FlagLoopback != 0 {
		t.Log().Debug().Msgf("skip arp announce on loopback interface %s", t.IpDev)
		return nil
	}
	t.Log().Info().Msgf("send gratuitous arp to announce %s over %s", t.ipaddr(), t.IpDev)
	return t.arpGratuitous()
}

func (t *T) start() error {
	ipnet := t.ipnet()
	if ipnet.Mask == nil {
		err := fmt.Errorf("ipnet definition error: %s/%s", t.ipaddr(), t.ipmask())
		t.Log().Error().Err(err).Msg("")
		return err
	}
	t.Log().Info().Msgf("add %s to %s", ipnet, t.IpDev)
	return netif.AddAddr(t.IpDev, ipnet)
}

func (t *T) stop() error {
	ipnet := t.ipnet()
	if ipnet.Mask == nil {
		err := fmt.Errorf("ipnet definition error: %s/%s", t.ipaddr(), t.ipmask())
		t.Log().Error().Err(err).Msg("")
		return err
	}
	t.Log().Info().Msgf("delete %s from %s", t.ipnet(), t.IpDev)
	return netif.DelAddr(t.IpDev, t.ipnet())
}
