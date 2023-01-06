package protocol

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/bwmarrin/discordgo"
	"github.com/rumblefrog/source-chat-relay/server/config"
	"github.com/rumblefrog/source-chat-relay/server/packet"
)

type IdentificationType uint8

const (
	IdentificationInvalid IdentificationType = iota
	IdentificationSteam
	IdentificationDiscord
	IdentificationTypeCount
)

type ChatMessage struct {
	BaseMessage

	IDType IdentificationType

	ID string

	Username string

	Message string
}

func ParseChatMessage(base BaseMessage, r *packet.PacketReader) (*ChatMessage, error) {
	m := &ChatMessage{}

	m.BaseMessage = base

	m.IDType = ParseIdentificationType(r.ReadUint8())

	var ok bool

	m.ID, ok = r.TryReadString()

	if !ok {
		return nil, ErrCannotReadString
	}

	m.Username, ok = r.TryReadString()

	if !ok {
		return nil, ErrCannotReadString
	}

	m.Message, ok = r.TryReadString()

	if !ok {
		return nil, ErrCannotReadString
	}

	return m, nil
}

func (m *ChatMessage) Type() MessageType {
	return MessageChat
}

func (m *ChatMessage) Content() string {
	return m.Message
}

func (m *ChatMessage) Marshal() []byte {
	var builder packet.PacketBuilder

	builder.WriteByte(byte(MessageChat))
	builder.WriteCString(m.BaseMessage.EntityName)

	builder.WriteByte(byte(m.IDType))
	builder.WriteCString(m.ID)
	builder.WriteCString(m.Username)
	builder.WriteCString(m.Message)

	return builder.Bytes()
}

func (m *ChatMessage) Plain() string {
	return strings.ReplaceAll(strings.ReplaceAll(config.Config.Messages.EventFormatSimplePlayerChat, "%username%", m.Username), "%message%", m.Message)
}

func (m *ChatMessage) Embed() *discordgo.MessageEmbed {
	// Generate Random colors by SteamID for Embed.
	// idColorBytes := []byte(m.ID)
	// Convert to an int with length of 6
	// color := int(binary.LittleEndian.Uint32(idColorBytes[len(idColorBytes)-6:])) / 10000

	steamURL := m.IDType.FormatURL(m.ID)

	steamProfileURL := SteamProfileScrape(steamURL)

	return &discordgo.MessageEmbed{
		Color:       16711784,
		Description: m.Message,
		Timestamp:   time.Now().Format(time.RFC3339),
		Author: &discordgo.MessageEmbedAuthor{
			IconURL: func() string {
				if steamProfileURL != "" {
					return steamProfileURL
				} else {
					// Question mark avatars from Steam.
					return "https://steamcdn-a.akamaihd.net/steamcommunity/public/images/avatars/fe/fef49e7fa7e1997310d705b2a6158ff8dc1cdfeb_full.jpg"
				}
			}(),
			Name: fmt.Sprintf("%s (%s)", m.Username, m.ID),
			URL:  m.IDType.FormatURL(m.ID),
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: m.BaseMessage.EntityName,
		},
	}
}

func SteamProfileScrape(steamURL string) string {
	// Request the HTML page
	var src string
	res, err := http.Get(steamURL)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("Failed to fetch data: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, _ := goquery.NewDocumentFromReader(res.Body)

	doc.Find(".playerAvatarAutoSizeInner img").Each(func(i int, s *goquery.Selection) {
		src, _ = s.Attr("src")

		if strings.Contains(src, "steamcommunity/public/images/items/") {
			src = ""
		}
	})

	return src
}

func ParseIdentificationType(t uint8) IdentificationType {
	if t >= uint8(IdentificationTypeCount) {
		return IdentificationInvalid
	}

	return IdentificationType(t)
}

func (i IdentificationType) FormatURL(id string) string {
	switch i {
	case IdentificationSteam:
		return "https://steamcommunity.com/profiles/" + id
	default:
		return ""
	}
}
