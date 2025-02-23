package object

import (
	"context"

	"opensvc.com/opensvc/core/actioncontext"
	"opensvc.com/opensvc/core/resource"
)

// SyncResync re-establishes the data synchronization
func (t *actor) SyncResync(ctx context.Context) error {
	ctx = actioncontext.WithProps(ctx, actioncontext.SyncResync)
	if err := t.validateAction(); err != nil {
		return err
	}
	t.setenv("sync_resync", false)
	unlock, err := t.lockAction(ctx)
	if err != nil {
		return err
	}
	defer unlock()
	return t.lockedSyncResync(ctx)
}

func (t *actor) lockedSyncResync(ctx context.Context) error {
	if err := t.masterSyncResync(ctx); err != nil {
		return err
	}
	if err := t.slaveSyncResync(ctx); err != nil {
		return err
	}
	return nil
}

func (t *actor) masterSyncResync(ctx context.Context) error {
	return t.action(ctx, func(ctx context.Context, r resource.Driver) error {
		return resource.Resync(ctx, r)
	})
}

func (t *actor) slaveSyncResync(ctx context.Context) error {
	return nil
}
