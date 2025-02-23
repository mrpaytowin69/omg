package commands

import (
	"opensvc.com/opensvc/core/nodeaction"
	"opensvc.com/opensvc/core/object"
)

type (
	CmdNodePrintConfig struct {
		OptsGlobal
		Eval        bool
		Impersonate string
	}
)

func (t *CmdNodePrintConfig) Run() error {
	return nodeaction.New(
		nodeaction.LocalFirst(),
		nodeaction.WithLocal(t.Local),
		nodeaction.WithRemoteNodes(t.NodeSelector),
		nodeaction.WithFormat(t.Format),
		nodeaction.WithColor(t.Color),
		nodeaction.WithServer(t.Server),
		nodeaction.WithRemoteAction("print config"),
		nodeaction.WithRemoteOptions(map[string]interface{}{
			"impersonate": t.Impersonate,
			"eval":        t.Eval,
		}),
		nodeaction.WithLocalRun(func() (interface{}, error) {
			n, err := object.NewNode()
			if err != nil {
				return nil, err
			}
			switch {
			case t.Eval:
				return n.EvalConfigAs(t.Impersonate)
			default:
				return n.PrintConfig()
			}
		}),
	).Do()
}
