package object

import (
	"context"

	"opensvc.com/opensvc/core/actioncontext"
	"opensvc.com/opensvc/core/resource"
)

// SetProvisioned starts the local instance of the object
func (t *actor) SetProvisioned(ctx context.Context) error {
	ctx = actioncontext.WithProps(ctx, actioncontext.SetProvisioned)
	if err := t.validateAction(); err != nil {
		return err
	}
	t.setenv("set provisioned", false)
	unlock, err := t.lockAction(ctx)
	if err != nil {
		return err
	}
	defer unlock()
	return t.lockedSetProvisioned(ctx)
}

func (t *actor) lockedSetProvisioned(ctx context.Context) error {
	if err := t.masterSetProvisioned(ctx); err != nil {
		return err
	}
	if err := t.slaveSetProvisioned(ctx); err != nil {
		return err
	}
	return nil
}

func (t *actor) masterSetProvisioned(ctx context.Context) error {
	return t.action(ctx, func(ctx context.Context, r resource.Driver) error {
		return resource.SetProvisioned(ctx, r)
	})
}

func (t *actor) slaveSetProvisioned(ctx context.Context) error {
	return nil
}
