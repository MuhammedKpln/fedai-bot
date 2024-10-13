package main

import (
	"fmt"
	C "muhammedkpln/fedai/core"
	M "muhammedkpln/fedai/modules"
	S "muhammedkpln/fedai/shared"

	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow/types"

	"go.mau.fi/whatsmeow/types/events"
)

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
	var textMessage = message.Message.GetConversation()

	// Catch only messages
	if message.Message.Conversation != nil || m != nil {
		var ci = m.GetContextInfo()

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

		textMessageWrapper := S.If(m != nil, *m.Text, textMessage)

		pl := CommandCatcher(textMessageWrapper)
		if pl != nil {
			fmt.Println(context.SenderJID, C.GetClient().Store.ID)
			// fmt.Println(bool(*pl.IsPublic))

			if pl.IsPublic == nil {
				if message.Info.Sender.ToNonAD() != C.GetClient().Store.ID.ToNonAD() {
					// Plugin is not allowed to be used of other users, return
					return
				}

			}
			fmt.Println("selam")
			go pl.CommandFn(&context)

		}
	}

}

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		go handleMessageEvent(v)

	case *events.Connected:
		C.AppLog().Infof("Connection established")

	}
}

func main() {
	C.LoadDotenv()
	M.LoadModules()
	C.AppLog().Infof("Modules loaded")
	C.EstablishConnection(eventHandler)

}
