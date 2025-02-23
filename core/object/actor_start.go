package object

import (
	"context"

	"opensvc.com/opensvc/core/actioncontext"
	"opensvc.com/opensvc/core/resource"
)

// Start starts the local instance of the object
func (t *actor) Start(ctx context.Context) error {
	ctx = actioncontext.WithProps(ctx, actioncontext.Start)
	if err := t.validateAction(); err != nil {
		return err
	}
	t.setenv("start", false)
	unlock, err := t.lockAction(ctx)
	if err != nil {
		return err
	}
	defer unlock()
	return t.lockedStart(ctx)
}

func (t *actor) lockedStart(ctx context.Context) error {
	if err := t.masterStart(ctx); err != nil {
		return err
	}
	if err := t.slaveStart(ctx); err != nil {
		return err
	}
	return nil
}

func (t *actor) masterStart(ctx context.Context) error {
	return t.action(ctx, func(ctx context.Context, r resource.Driver) error {
		t.log.Debug().Str("rid", r.RID()).Msg("start resource")
		return resource.Start(ctx, r)
	})
}

func (t *actor) slaveStart(ctx context.Context) error {
	return nil
}
