package bot

import (
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/rumblefrog/source-chat-relay/server/config"
	"github.com/rumblefrog/source-chat-relay/server/relay"
)

func aboutCommand(ctx *exrouter.Context) {
	ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "Fishy!",
			URL:     "https://github.com/rumblefrog",
			IconURL: "https://avatars2.githubusercontent.com/u/6960234?s=32",
		},
		Color: 0x3395D6,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "SCR Version",
				Value: config.SCRVER,
			},
			{
				Name:  "Repository",
				Value: "https://github.com/rumblefrog/source-chat-relay/",
			},
		},
	})

	ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "Suganda",
			URL:     "https://github.com/suganda8",
			IconURL: "https://avatars.githubusercontent.com/u/58646818?s=32",
		},
		Color: 16711784,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Forked SCR Version",
				Value: config.SCRVER,
			},
			{
				Name:  "New Feature",
				Value: "- Embed now has user's profile image\n- Decode Discord message that contain emojis to its alias \\:heart\\: \\:heart\\: \\:heart\\:",
			},
			{
				Name:  "Repository",
				Value: "https://github.com/suganda8/mep-source-chat-relay",
			},
		},
	})

	ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, &discordgo.MessageEmbed{
		Title: "Traffics",
		Color: 16711784,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Incoming Traffic",
				Value: relay.Instance.Statistics.Incoming.String(),
			},
			{
				Name:  "Outgoing Traffic",
				Value: relay.Instance.Statistics.Outgoing.String(),
			},
		},
	})
}
