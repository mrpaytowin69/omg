package object

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"opensvc.com/opensvc/core/keyop"
	"opensvc.com/opensvc/core/rawconfig"
	"opensvc.com/opensvc/core/xconfig"
	"opensvc.com/opensvc/util/hostname"
	"opensvc.com/opensvc/util/key"
)

func (t Node) Log() *zerolog.Logger {
	return &t.log
}

func (t Node) Exists() bool {
	return true
}

func (t *Node) ConfigFile() string {
	return rawconfig.NodeConfigFile()
}

func (t *Node) ClusterConfigFile() string {
	return rawconfig.ClusterConfigFile()
}

func (t *Node) loadConfig() error {
	nodeConfigFile := t.ConfigFile()
	if config, err := xconfig.NewObject(nodeConfigFile, nodeConfigFile); err != nil {
		return err
	} else {
		t.config = config
		t.config.Referrer = t
	}
	clusterConfigFile := t.ClusterConfigFile()
	if config, err := xconfig.NewObject(clusterConfigFile, clusterConfigFile, nodeConfigFile); err != nil {
		return err
	} else {
		t.mergedConfig = config
		t.mergedConfig.Referrer = t
	}
	return nil
}

func (t Node) Config() *xconfig.T {
	return t.config
}

func (t Node) MergedConfig() *xconfig.T {
	return t.mergedConfig
}

func (t Node) ID() uuid.UUID {
	if t.id != uuid.Nil {
		return t.id
	}
	idKey := key.Parse("id")
	if idStr := t.config.GetString(idKey); idStr != "" {
		if id, err := uuid.Parse(idStr); err == nil {
			t.id = id
			return t.id
		}
	}
	t.id = uuid.New()
	op := keyop.T{
		Key:   key.Parse("id"),
		Op:    keyop.Set,
		Value: t.id.String(),
	}
	_ = t.config.Set(op)
	if err := t.config.Commit(); err != nil {
		t.log.Error().Err(err).Msg("")
	}
	return t.id
}

func (t Node) Env() string {
	k := key.Parse("env")
	if s := t.config.GetString(k); s != "" {
		return s
	}
	return rawconfig.NodeSection().Env
}

func (t Node) App() string {
	k := key.Parse("app")
	return t.config.GetString(k)
}

func (t Node) Dereference(ref string) (string, error) {
	switch ref {
	case "id":
		return t.ID().String(), nil
	case "name", "nodename":
		return hostname.Hostname(), nil
	case "short_name", "short_nodename":
		return strings.SplitN(hostname.Hostname(), ".", 1)[0], nil
	case "dnsuxsock":
		return rawconfig.DNSUDSFile(), nil
	case "dnsuxsockd":
		return rawconfig.DNSUDSDir(), nil
	}
	switch {
	case strings.HasPrefix(ref, "safe://"):
		return ref, fmt.Errorf("TODO")
	}
	return ref, fmt.Errorf("unknown reference: %s", ref)
}

func (t Node) PostCommit() error {
	return nil
}

func (t Node) Nodes() []string {
	k := key.T{Section: "cluster", Option: "nodes"}
	nodes := t.MergedConfig().GetStrings(k)
	if len(nodes) == 0 {
		return []string{hostname.Hostname()}
	}
	return nodes
}

func (t Node) DRPNodes() []string {
	return []string{}
}

func (t Node) EncapNodes() []string {
	return []string{}
}

func (t *Node) Nameservers() ([]string, error) {
	dns, err := t.MergedConfig().Eval(key.T{Section: "cluster", Option: "dns"})
	return dns.([]string), err
}

func (t *Node) CNIConfig() (string, error) {
	if s, err := t.MergedConfig().Eval(key.T{Section: "cni", Option: "config"}); err != nil {
		return "", err
	} else {
		return s.(string), nil
	}
}

func (t *Node) CNIPlugins() (string, error) {
	if s, err := t.MergedConfig().Eval(key.T{Section: "cni", Option: "plugins"}); err != nil {
		return "", err
	} else {
		return s.(string), nil
	}
}

func (t *Node) Labels() map[string]string {
	return t.config.SectionMap("labels")
}
