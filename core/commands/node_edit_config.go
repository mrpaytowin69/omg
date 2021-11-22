package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"opensvc.com/opensvc/core/flag"
	"opensvc.com/opensvc/core/object"
)

type (
	// CmdObjectEditConfig is the cobra flag set of the print config command.
	NodeEditConfig struct {
		Global     object.OptsGlobal
		EditConfig object.OptsEditConfig
	}
)

// Init configures a cobra command and adds it to the parent command.
func (t *NodeEditConfig) Init(parent *cobra.Command) {
	cmd := t.cmd()
	parent.AddCommand(cmd)
	flag.Install(cmd, t)
}

func (t *NodeEditConfig) cmd() *cobra.Command {
	return &cobra.Command{
		Use:     "config",
		Short:   "edit the node configuration",
		Aliases: []string{"confi", "conf", "con", "co", "c", "cf", "cfg"},
		Run: func(_ *cobra.Command, _ []string) {
			t.run()
		},
	}
}

func (t *NodeEditConfig) run() {
	if err := object.NewNode().EditConfig(t.EditConfig); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
