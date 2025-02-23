package commands

import (
	"opensvc.com/opensvc/core/nodeaction"
	"opensvc.com/opensvc/core/object"
)

type (
	CmdNodePrintCapabilities struct {
		OptsGlobal
	}
)

func (t *CmdNodePrintCapabilities) Run() error {
	return nodeaction.New(
		nodeaction.WithFormat(t.Format),
		nodeaction.WithColor(t.Color),
		nodeaction.WithServer(t.Server),

		nodeaction.WithRemoteNodes(t.NodeSelector),
		nodeaction.WithRemoteAction("node print capabilities"),
		nodeaction.WithRemoteOptions(map[string]interface{}{
			"format": t.Format,
		}),

		nodeaction.WithLocal(t.Local),
		nodeaction.WithLocalRun(func() (interface{}, error) {
			n, err := object.NewNode()
			if err != nil {
				return nil, err
			}
			return n.PrintCapabilities()
		}),
	).Do()
}
