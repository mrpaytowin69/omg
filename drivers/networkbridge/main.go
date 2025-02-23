package networkbridge

import (
	"opensvc.com/opensvc/core/driver"
	"opensvc.com/opensvc/core/network"
)

type (
	T struct {
		network.T
	}
)

var (
	drvID = driver.NewID(driver.GroupNetwork, "bridge")
)

func init() {
	driver.Register(drvID, NewNetworker)
}

func NewNetworker() network.Networker {
	t := New()
	var i interface{} = t
	return i.(network.Networker)
}

func New() *T {
	t := T{}
	return &t
}

func (t T) Usage() (network.StatusUsage, error) {
	usage := network.StatusUsage{}
	return usage, nil
}

func (t T) brName() string {
	return "obr_" + t.Name()
}

func (t T) BackendDevName() string {
	return t.brName()
}

// CNIConfigData returns a cni network configuration, like
// {
//   "bridge": "cni0",
//   "cniVersion": "0.3.0",
//   "ipMasq": true,
//   "name": "mynet",
//   "ipam": {
//     "routes": [
//       {"dst": "0.0.0.0/0"}
//     ],
//     "subnet": "10.22.0.0/16",
//     "type": "host-local"
//   },
//   "isGateway": true,
//   "type": "bridge"
// }
func (t T) CNIConfigData() (interface{}, error) {
	m := map[string]interface{}{
		"cniVersion": network.CNIVersion,
		"name":       t.Name(),
		"type":       "bridge",
		"bridge":     t.brName(),
		"isGateway":  true,
		"ipMasq":     true,
		"ipam": map[string]interface{}{
			"type": "host-local",
			"routes": []map[string]interface{}{
				{"dst": "0.0.0.0/0"},
			},
			"subnet": t.Network(),
		},
	}
	return m, nil
}
