package commands

import (
	"opensvc.com/opensvc/core/client"
	"opensvc.com/opensvc/core/nodeselector"
	"opensvc.com/opensvc/core/output"
	"opensvc.com/opensvc/core/rawconfig"
)

type (
	CmdNodeLs struct {
		OptsGlobal
	}
)

func (t *CmdNodeLs) Run() error {
	var (
		c        *client.T
		err      error
		selector string
	)
	if !t.Local {
		if c, err = client.New(client.WithURL(t.Server)); err != nil {
			return err
		}
	}
	if t.NodeSelector == "" {
		selector = "*"
	} else {
		selector = t.NodeSelector
	}
	nodes := nodeselector.New(
		selector,
		nodeselector.WithLocal(t.Local),
		nodeselector.WithServer(t.Server),
		nodeselector.WithClient(c),
	).Expand()
	output.Renderer{
		Format: t.Format,
		Color:  t.Color,
		Data:   nodes,
		HumanRenderer: func() string {
			s := ""
			for _, e := range nodes {
				s += e + "\n"
			}
			return s
		},
		Colorize: rawconfig.Colorize,
	}.Print()
	return nil
}
