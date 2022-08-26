package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/spf13/cobra"
	"opensvc.com/opensvc/core/client"
	"opensvc.com/opensvc/core/flag"
	"opensvc.com/opensvc/core/nodeselector"
	"opensvc.com/opensvc/core/rawconfig"
	"opensvc.com/opensvc/core/slog"
	"opensvc.com/opensvc/util/render"
)

type (
	// NodeLogs is the cobra flag set of the logs command.
	NodeLogs struct {
		OptsGlobal
		Follow bool   `flag:"logs-follow"`
		SID    string `flag:"logs-sid"`
	}
)

// Init configures a cobra command and adds it to the parent command.
func (t *NodeLogs) Init(parent *cobra.Command) {
	cmd := t.cmd()
	parent.AddCommand(cmd)
	flag.Install(cmd, t)
}

func (t *NodeLogs) cmd() *cobra.Command {
	return &cobra.Command{
		Use:     "logs",
		Aliases: []string{"logs", "log", "lo"},
		Short:   "filter and format logs",
		Run: func(cmd *cobra.Command, args []string) {
			t.run()
		},
	}
}

func (t *NodeLogs) backlog(node string) (slog.Events, error) {
	c, err := client.New(
		client.WithURL(node),
		client.WithPassword(rawconfig.ClusterSection().Secret),
	)
	if err != nil {
		return nil, err
	}
	req := c.NewGetNodeBacklog().SetFilters(t.Filters())
	b, err := req.Do()
	if err != nil {
		return nil, err
	}
	events := make(slog.Events, 0)
	if err := json.Unmarshal(b, &events); err != nil {
		return nil, err
	}
	return events, nil
}

func (t *NodeLogs) stream(node string) {
	c, err := client.New(
		client.WithURL(node),
		client.WithPassword(rawconfig.ClusterSection().Secret),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	streamer := c.NewGetNodeLog().SetFilters(t.Filters())
	events, err := streamer.Do()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	for event := range events {
		event.Render(t.Format)
	}
}

func (t *NodeLogs) remote() error {
	sel := nodeselector.New(
		t.NodeSelector,
		nodeselector.WithServer(t.Server),
	)
	nodes := sel.Expand()
	if len(nodes) == 0 {
		return errors.New("no nodes to fetch logs from")
	}
	events := make(slog.Events, 0)
	for _, node := range nodes {
		if more, err := t.backlog(node); err != nil {
			fmt.Fprintln(os.Stderr, "backlog fetch error:", err)
		} else {
			events = append(events, more...)
		}
	}
	events.Sort()
	events.Render(t.Format)
	if !t.Follow {
		return nil
	}
	var wg sync.WaitGroup
	wg.Add(len(nodes))
	for _, node := range nodes {
		go func() {
			defer wg.Done()
			t.stream(node)
		}()
	}
	wg.Wait()
	return nil
}

func (t NodeLogs) Filters() map[string]interface{} {
	filters := make(map[string]interface{})
	if t.SID != "" {
		filters["sid"] = t.SID
	}
	return filters
}

func (t *NodeLogs) local() error {
	filters := t.Filters()
	if events, err := slog.GetEventsFromNode(filters); err == nil {
		events.Render(t.Format)
	} else {
		return err
	}
	if t.Follow {
		stream, err := slog.GetEventStreamFromNode(filters)
		if err != nil {
			return err
		}
		for event := range stream.Events() {
			event.Render(t.Format)
		}
	}
	return nil
}

func (t *NodeLogs) run() {
	var err error
	render.SetColor(t.Color)
	if t.NodeSelector == "" {
		t.NodeSelector = "*"
	}
	if t.Local {
		err = t.local()
	} else {
		err = t.remote()
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
