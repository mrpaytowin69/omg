package commands

import (
	"opensvc.com/opensvc/core/nodeaction"
	"opensvc.com/opensvc/core/object"
)

type (
	CmdNodeDelete struct {
		OptsGlobal
		RID string
	}
)

func (t *CmdNodeDelete) Run() error {
	return nodeaction.New(
		nodeaction.LocalFirst(),
		nodeaction.WithLocal(t.Local),
		nodeaction.WithRemoteNodes(t.NodeSelector),
		nodeaction.WithFormat(t.Format),
		nodeaction.WithColor(t.Color),
		nodeaction.WithServer(t.Server),
		nodeaction.WithRemoteAction("delete"),
		nodeaction.WithRemoteOptions(map[string]interface{}{
			"rid": t.RID,
		}),
		nodeaction.WithLocalRun(func() (interface{}, error) {
			n, err := object.NewNode()
			if err != nil {
				return nil, err
			}
			return nil, n.DeleteSection(t.RID)
		}),
	).Do()
}
