package commands

import (
	"github.com/spf13/cobra"
	"opensvc.com/opensvc/core/entrypoints/nodeaction"
	"opensvc.com/opensvc/core/flag"
	"opensvc.com/opensvc/core/object"
)

type (
	// NodePrintConfig is the cobra flag set of the start command.
	NodePrintConfig struct {
		object.OptsPrintConfig
	}
)

// Init configures a cobra command and adds it to the parent command.
func (t *NodePrintConfig) Init(parent *cobra.Command) {
	cmd := t.cmd()
	parent.AddCommand(cmd)
	flag.Install(cmd, &t.OptsPrintConfig)
}

func (t *NodePrintConfig) cmd() *cobra.Command {
	return &cobra.Command{
		Use:     "config",
		Short:   "get a configuration key value",
		Aliases: []string{"confi", "conf", "con", "co", "c", "cf", "cfg"},
		Run: func(_ *cobra.Command, _ []string) {
			t.run()
		},
	}
}

func (t *NodePrintConfig) run() {
	nodeaction.New(
		nodeaction.LocalFirst(),
		nodeaction.WithLocal(t.Global.Local),
		nodeaction.WithRemoteNodes(t.Global.NodeSelector),
		nodeaction.WithFormat(t.Global.Format),
		nodeaction.WithColor(t.Global.Color),
		nodeaction.WithServer(t.Global.Server),
		nodeaction.WithRemoteAction("print config"),
		nodeaction.WithRemoteOptions(map[string]interface{}{
			"impersonate": t.Impersonate,
			"eval":        t.Eval,
		}),
		nodeaction.WithLocalRun(func() (interface{}, error) {
			return object.NewNode().PrintConfig(t.OptsPrintConfig)
		}),
	).Do()
}
