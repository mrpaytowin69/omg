package commands

import (
	"context"

	"opensvc.com/opensvc/core/actioncontext"
	"opensvc.com/opensvc/core/object"
	"opensvc.com/opensvc/core/objectaction"
	"opensvc.com/opensvc/core/path"
)

type (
	CmdObjectStatus struct {
		OptsGlobal
		OptsLock
		Refresh bool
	}
)

func (t *CmdObjectStatus) Run(selector, kind string) error {
	mergedSelector := mergeSelector(selector, t.ObjectSelector, kind, "")
	return objectaction.New(
		objectaction.LocalFirst(),
		objectaction.WithObjectSelector(mergedSelector),
		objectaction.WithLocal(t.Local),
		objectaction.WithFormat(t.Format),
		objectaction.WithColor(t.Color),
		objectaction.WithRemoteNodes(t.NodeSelector),
		objectaction.WithRemoteAction("status"),
		objectaction.WithLocalRun(func(p path.T) (interface{}, error) {
			o, err := object.NewCore(p)
			if err != nil {
				return nil, err
			}
			ctx := context.Background()
			ctx = actioncontext.WithLockDisabled(ctx, t.Disable)
			ctx = actioncontext.WithLockTimeout(ctx, t.Timeout)
			if t.Refresh {
				return o.FreshStatus(ctx)
			} else {
				return o.Status(ctx)
			}
		}),
	).Do()
}
