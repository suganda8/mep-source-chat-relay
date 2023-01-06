package entity

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

func (entity *Entity) GetDisplayName() string {
	if len(entity.DisplayName) == 0 {
		return "N/A"
	}

	return entity.DisplayName
}

func (entity *Entity) Embed() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color: 0xE1C15C,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Entity ID :key:",
				Value: entity.ID,
			},
			{
				Name:  "Display Name :page_facing_up:",
				Value: entity.GetDisplayName(),
			},
			{
				Name:  "Receive Channels :inbox_tray:",
				Value: HumanizeChannelString(entity.ReceiveChannels),
			},
			{
				Name:  "Send Channels :outbox_tray:",
				Value: HumanizeChannelString(entity.SendChannels),
			},
			{
				Name:  "Created At",
				Value: entity.CreatedAt.Format(time.RFC1123),
			},
		},
	}
}
