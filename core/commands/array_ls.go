package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"opensvc.com/opensvc/core/client"
	"opensvc.com/opensvc/core/clientcontext"
	"opensvc.com/opensvc/core/flag"
	"opensvc.com/opensvc/core/object"
	"opensvc.com/opensvc/core/output"
	"opensvc.com/opensvc/core/rawconfig"
)

type (
	// ArrayLs is the cobra flag set of the command.
	ArrayLs struct {
		OptsGlobal
	}
)

// Init configures a cobra command and adds it to the parent command.
func (t *ArrayLs) Init(parent *cobra.Command) {
	cmd := t.cmd()
	parent.AddCommand(cmd)
	flag.Install(cmd, t)
}

func (t *ArrayLs) cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ls",
		Short: "list the cluster pools",
		Run: func(_ *cobra.Command, _ []string) {
			t.run()
		},
	}
}

func (t *ArrayLs) run() {
	var data []string
	if t.Local || !clientcontext.IsSet() {
		data = t.extractLocal()
	} else {
		data = t.extractDaemon()
	}
	output.Renderer{
		Format: t.Format,
		Color:  t.Color,
		Data:   data,
		HumanRenderer: func() string {
			s := ""
			for _, e := range data {
				s += e + "\n"
			}
			return s
		},
		Colorize: rawconfig.Colorize,
	}.Print()
}

func (t *ArrayLs) extractLocal() []string {
	n, err := object.NewNode()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return n.ListArrays()
}

func (t *ArrayLs) extractDaemon() []string {
	var (
		c   *client.T
		err error
	)
	if c, err = client.New(client.WithURL(t.Server)); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	panic("TODO")
	fmt.Println(c)
	return []string{}
}
