package object

import (
	"context"

	"opensvc.com/opensvc/core/actioncontext"
	"opensvc.com/opensvc/core/resource"
)

// Unprovision stops and frees the local instance of the object
func (t *actor) Unprovision(ctx context.Context) error {
	ctx = actioncontext.WithProps(ctx, actioncontext.Unprovision)
	if err := t.validateAction(); err != nil {
		return err
	}
	t.setenv("unprovision", false)
	unlock, err := t.lockAction(ctx)
	if err != nil {
		return err
	}
	defer unlock()
	return t.lockedUnprovision(ctx)
}

func (t *actor) lockedUnprovision(ctx context.Context) error {
	if err := t.slaveUnprovision(ctx); err != nil {
		return err
	}
	if err := t.masterUnprovision(ctx); err != nil {
		return err
	}
	return nil
}

func (t *actor) masterUnprovision(ctx context.Context) error {
	return t.action(ctx, func(ctx context.Context, r resource.Driver) error {
		t.log.Debug().Str("rid", r.RID()).Msg("unprovision resource")
		leader := actioncontext.IsLeader(ctx)
		return resource.Unprovision(ctx, r, leader)
	})
}

func (t *actor) slaveUnprovision(ctx context.Context) error {
	return nil
}
