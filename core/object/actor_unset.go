package object

import (
	"context"

	"opensvc.com/opensvc/core/actioncontext"
	"opensvc.com/opensvc/core/statusbus"
	"opensvc.com/opensvc/util/key"
)

// Unset gets a keyword value
func (t *actor) Unset(ctx context.Context, kws ...key.T) error {
	ctx = actioncontext.WithProps(ctx, actioncontext.Unset)
	ctx, stop := statusbus.WithContext(ctx, t.path)
	defer stop()
	defer t.postActionStatusEval(ctx)
	unlock, err := t.lockAction(ctx)
	if err != nil {
		return err
	}
	defer unlock()
	return unsetKeys(t.config, kws...)
}
