package main

import (
	C "muhammedkpln/fedai/core"
	M "muhammedkpln/fedai/modules"
	S "muhammedkpln/fedai/shared"
	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"

	"go.mau.fi/whatsmeow/types/events"
)

var appLog = waLog.Stdout("APP", "INFO", true)

func CommandCatcher(message string) *S.Plugin {
	var pl *S.Plugin

	for _, plugin := range M.LoadedPlugins {
		var matches = plugin.CommandRegex.MatchString(message)

		if matches {
			pl = &plugin
			break
		}
	}

	if pl != nil {
		return pl
	}

	return pl
}

func handleMessageEvent(message *events.Message) {
	// set as unavailable to not see(?) the message
	go C.GetClient().SendChatPresence(message.Info.Chat, "unavailable", types.ChatPresenceMediaText)
	var m = message.Message.ExtendedTextMessage
	var textMessage = m.GetText()

	// Catch only messages
	if m != nil {
		var ci = m.ContextInfo

		var context S.PluginRunOptions = S.PluginRunOptions{
			IsQuoted:  false,
			Message:   message,
			ChatJID:   message.Info.Chat,
			SenderJID: message.Info.Sender,
			StanzaID:  message.Info.ID,
		}

		// Checks if message has context, which can mean that it has a quoted message.
		if ci != nil {
			quotedMessage := ci.QuotedMessage
			if quotedMessage != nil {
				context = S.PluginRunOptions{
					IsQuoted:      true,
					QuotedMessage: quotedMessage,
					Message:       message,
					ChatJID:       message.Info.Chat,
					SenderJID:     message.Info.Sender,
					StanzaID:      message.Info.ID,
				}
			}

		}

		pl := CommandCatcher(textMessage)
		if pl != nil {
			go pl.CommandFn(&context)
		}
	}

}

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		handleMessageEvent(v)

	case *events.Connected:
		appLog.Infof("Connection established")

	}
}

func main() {
	M.LoadModules()
	appLog.Infof("Modules loaded")
	C.EstablishConnection(eventHandler)

}
