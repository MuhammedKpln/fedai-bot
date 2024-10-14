package main

import (
	"context"
	"fmt"
	Cx "muhammedkpln/fedai/context"
	C "muhammedkpln/fedai/core"
	M "muhammedkpln/fedai/modules"
	S "muhammedkpln/fedai/shared"
	"regexp"
	"time"

	"go.mau.fi/whatsmeow/proto/waE2E"
)

var Plugin S.Plugin = S.Plugin{
	Name:         "Help",
	CommandRegex: regexp.MustCompile(".help"),
	CommandInfo:  "Help",
	CommandFn:    Run,
}

func Run(message *S.PluginRunOptions, payload S.RegexpMatches) {
	var helpString string
	client := C.GetClient()
	for _, plugin := range M.LoadedPlugins {
		helpString += fmt.Sprintf("_*%s*_ \n *About*: %s \n Command: *%s* \n\n", plugin.Name, plugin.CommandInfo, plugin.CommandRegex.String())
	}

	go Cx.EditMessage(Cx.InfoMessage("Kendi sohbetine bak!"), message)
	go client.SendMessage(context.Background(), message.Message.Info.Sender.ToNonAD(), &waE2E.Message{
		Conversation: &helpString,
	})
	time.Sleep(3 * time.Second)
	go client.SendMessage(context.Background(), message.Message.Info.Chat, client.BuildRevoke(message.ChatJID, message.SenderJID, message.Message.Info.ID))

}
