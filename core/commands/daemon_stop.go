package commands

import "opensvc.com/opensvc/daemon/daemoncli"

type (
	CmdDaemonStop struct {
		OptsGlobal
	}
)

func (t *CmdDaemonStop) Run() error {
	cli, err := newClient(t.Server)
	if err != nil {
		return err
	}
	daemoncli.LockFuncExit("daemon stop", daemoncli.New(cli).Stop)
	return nil
}
