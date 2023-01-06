package bot

import (
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/rumblefrog/source-chat-relay/server/protocol"
	"github.com/rumblefrog/source-chat-relay/server/relay"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func eventCommand(ctx *exrouter.Context) {
	if len(ctx.Args) < 2 {
		ctx.Reply("Missing arguments")

		return
	}

	channel, err := ctx.Channel(ctx.Msg.ChannelID)

	if err != nil {
		ctx.Reply("Unable to fetch channel")

		return
	}

	caser := cases.Title(language.Und, cases.NoLower)

	message := &protocol.EventMessage{
		BaseMessage: protocol.BaseMessage{
			Type:       protocol.MessageChat,
			SenderID:   ctx.Msg.ChannelID,
			EntityName: caser.String(channel.Name),
		},
		Event: ctx.Args.Get(0),
		Data:  ctx.Args.Get(1),
	}

	relay.Instance.Router <- message
}
