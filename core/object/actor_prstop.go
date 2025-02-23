package object

import (
	"context"

	"opensvc.com/opensvc/core/actioncontext"
	"opensvc.com/opensvc/core/resource"
)

// PRStop stops the exclusive write access to devices of the local instance of the object
func (t *actor) PRStop(ctx context.Context) error {
	ctx = actioncontext.WithProps(ctx, actioncontext.Stop)
	if err := t.validateAction(); err != nil {
		return err
	}
	t.setenv("start", false)
	unlock, err := t.lockAction(ctx)
	if err != nil {
		return err
	}
	defer unlock()
	return t.lockedPRStop(ctx)
}

func (t *actor) lockedPRStop(ctx context.Context) error {
	if err := t.masterPRStop(ctx); err != nil {
		return err
	}
	if err := t.slavePRStop(ctx); err != nil {
		return err
	}
	return nil
}

func (t *actor) masterPRStop(ctx context.Context) error {
	return t.action(ctx, func(ctx context.Context, r resource.Driver) error {
		t.log.Debug().Str("rid", r.RID()).Msg("start resource")
		return resource.PRStop(ctx, r)
	})
}

func (t *actor) slavePRStop(ctx context.Context) error {
	return nil
}
