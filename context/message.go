package context

import (
	"context"

	C "muhammedkpln/fedai/core"
	S "muhammedkpln/fedai/shared"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

func SendQuotedMessage(chatJID types.JID, text string, senderJID string, infoID *string, quotedMessage *waE2E.Message) {
	client := C.GetClient()

	go client.SendMessage(context.Background(), chatJID, &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: proto.String(text),
			ContextInfo: &waE2E.ContextInfo{
				StanzaID:      infoID,
				QuotedMessage: quotedMessage,
				Participant:   proto.String(senderJID),
			},
		},
	})
}

func SendMessage(chatJID types.JID, text string, senderJID string, infoID *string, quotedMessage *waE2E.Message) {
	client := C.GetClient()

	go client.SendMessage(context.Background(), chatJID, &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: proto.String(text),
			ContextInfo: &waE2E.ContextInfo{
				StanzaID:    infoID,
				Participant: proto.String(senderJID),
			},
		},
	})
}

func EditMessage(text string, message *S.PluginRunOptions) {
	client := C.GetClient()
	chatJID := message.ChatJID
	infoID := message.StanzaID

	go client.SendMessage(context.Background(), chatJID, client.BuildEdit(chatJID, infoID, &waE2E.Message{
		Conversation: &text,
	}))
}

func SuccessMessage(msg string) string {
	return "✅ *FEDAI BOT*:  ```" + msg + "```"
}
func ErrorMessage(msg string) string {
	return "🛑 *FEDAI BOT*:  ```" + msg + "```"
}
func InfoMessage(msg string) string {
	return "⏺️ *FEDAI BOT*:  ```" + msg + "```"
}
