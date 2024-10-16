package shared

import (
	"regexp"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

type RegexpMatches struct {
	Match   string
	Action  *string
	Payload *string
}

type Plugin struct {
	Name         string
	CommandRegex *regexp.Regexp
	CommandInfo  string
	CommandFn    func(*PluginRunOptions, RegexpMatches)
	IsPublic     *bool
}

type PluginRunOptions struct {
	IsQuoted      bool
	QuotedMessage *waE2E.Message
	Message       *events.Message
	ChatJID       types.JID
	SenderJID     types.JID
	StanzaID      string
}
