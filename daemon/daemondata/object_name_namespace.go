package daemondata

import (
	"context"

	"opensvc.com/opensvc/core/path"
)

type opGetServiceNames struct {
	services chan<- []string
}

// GetServiceNames returns the cluster wide list of path.T.String() parsed from the cluster dataset, in
// cluster.node[*].instance[*].config
func (t T) GetServiceNames() []string {
	services := make(chan []string)
	t.cmdC <- opGetServiceNames{
		services: services,
	}
	return <-services
}

// GetServicePaths returns the cluster wide path.L parsed from the cluster dataset, in
// cluster.node[*].Services.Config[*]
func (t T) GetServicePaths() path.L {
	l := t.GetServiceNames()
	paths, _ := path.ParseList(l...)
	return paths
}

func (t T) GetNamespaces() []string {
	return t.GetServicePaths().Namespaces()
}

func (o opGetServiceNames) call(ctx context.Context, d *data) {
	paths := make(map[string]bool)
	for node := range d.pending.Cluster.Node {
		for s, inst := range d.pending.Cluster.Node[node].Instance {
			if inst.Config != nil {
				paths[s] = true
			}
		}
	}
	services := make([]string, 0)
	for s := range paths {
		services = append(services, s)
	}
	select {
	case <-ctx.Done():
	case o.services <- services:
	}
}

type (
	TInstanceId = string
)

func InstanceId(p path.T, node string) TInstanceId {
	return TInstanceId(node + ":" + p.String())
}
