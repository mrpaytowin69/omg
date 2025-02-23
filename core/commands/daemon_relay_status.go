package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"opensvc.com/opensvc/core/client"
	"opensvc.com/opensvc/core/object"
	"opensvc.com/opensvc/core/output"
	"opensvc.com/opensvc/core/path"
	"opensvc.com/opensvc/core/rawconfig"
	"opensvc.com/opensvc/core/xconfig"
	"opensvc.com/opensvc/daemon/daemonapi"
	"opensvc.com/opensvc/util/key"
	"opensvc.com/opensvc/util/render/tree"
)

type (
	CmdDaemonRelayStatus struct {
		OptsGlobal
		Relays string
	}
	relayMessage struct {
		daemonapi.RelayMessage
		Relay string `json:"relay"`
	}
	relayMessages []relayMessage
)

func (t *CmdDaemonRelayStatus) Run() error {
	messages := make(relayMessages, 0)
	relayMap := make(map[string]any)
	if t.Relays != "" {
		for _, s := range strings.Split(t.Relays, ",") {
			relayMap[s] = nil
		}
	}
	node, err := object.NewNode()
	if err != nil {
		return err
	}
	config := node.MergedConfig()
	for _, section := range config.SectionStrings() {
		if !strings.HasPrefix(section, "hb#") {
			continue
		}
		hbType := config.Get(key.New(section, "type"))
		if hbType != "relay" {
			continue
		}
		hbRelay := config.GetString(key.New(section, "relay"))
		if len(relayMap) > 0 {
			// some relay filtering is on
			if _, ok := relayMap[hbRelay]; !ok {
				// filtered out
				continue
			}
		}
		insecure := config.GetBool(key.New(section, "insecure"))
		username := config.GetString(key.New(section, "username"))
		password, err := configSectionPassword(config, section)
		if err != nil {
			return err
		}
		cli, err := client.New(
			client.WithURL(hbRelay),
			client.WithUsername(username),
			client.WithPassword(password),
			client.WithInsecureSkipVerify(insecure),
		)
		if err != nil {
			return err
		}
		req := cli.NewGetRelayMessage()
		b, err := req.Do()
		if err != nil {
			return err
		}
		var data daemonapi.RelayMessages
		if err := json.Unmarshal(b, &data); err != nil {
			return err
		}
		for _, message := range data.Messages {
			messages = append(messages, relayMessage{
				Relay:        hbRelay,
				RelayMessage: message,
			})
		}
	}
	output.Renderer{
		Format:   t.Format,
		Color:    t.Color,
		Data:     messages,
		Colorize: rawconfig.Colorize,
		HumanRenderer: func() string {
			return messages.Render()
		},
	}.Print()
	return nil
}

func (t relayMessages) Render() string {
	tree := tree.New()
	tree.AddColumn().AddText("Relay").SetColor(rawconfig.Color.Bold)
	tree.AddColumn().AddText("ClusterId").SetColor(rawconfig.Color.Bold)
	tree.AddColumn().AddText("ClusterName").SetColor(rawconfig.Color.Bold)
	tree.AddColumn().AddText("NodeName").SetColor(rawconfig.Color.Bold)
	tree.AddColumn().AddText("NodeAddr").SetColor(rawconfig.Color.Bold)
	tree.AddColumn().AddText("UpdatedAt").SetColor(rawconfig.Color.Bold)
	tree.AddColumn().AddText("MessageLength").SetColor(rawconfig.Color.Bold)
	for _, e := range t {
		n := tree.AddNode()
		n.AddColumn().AddText(e.Relay).SetColor(rawconfig.Color.Primary)
		n.AddColumn().AddText(e.ClusterId)
		n.AddColumn().AddText(e.ClusterName).SetColor(rawconfig.Color.Primary)
		n.AddColumn().AddText(e.Nodename).SetColor(rawconfig.Color.Primary)
		n.AddColumn().AddText(e.Addr)
		n.AddColumn().AddText(fmt.Sprint(e.Updated))
		n.AddColumn().AddText(fmt.Sprint(len(e.Msg)))
	}
	return tree.Render()
}

func configSectionPasswordSec(config *xconfig.T, section string) (object.Sec, error) {
	s := config.GetString(key.New(section, "password"))
	secPath, err := path.Parse(s)
	if err != nil {
		return nil, err
	}
	return object.NewSec(secPath, object.WithVolatile(true))
}

func configSectionPassword(config *xconfig.T, section string) (string, error) {
	sec, err := configSectionPasswordSec(config, section)
	if err != nil {
		return "", err
	}
	b, err := sec.DecodeKey("password")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
