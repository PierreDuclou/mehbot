package util

import (
	"log"
	"mehbot/config"

	"github.com/bwmarrin/discordgo"
)

type Message struct {
	Type         int
	Content      string
	MessageEmbed *discordgo.MessageEmbed
}

func NewMessage(messageType int) *Message {
	return &Message{Type: messageType, MessageEmbed: &discordgo.MessageEmbed{}}
}

func (m *Message) WithTitle(title string) *Message {
	m.MessageEmbed.Title = title
	return m
}

func (m *Message) WithFields(fields ...*discordgo.MessageEmbedField) *Message {
	for _, field := range fields {
		m.MessageEmbed.Fields = append(m.MessageEmbed.Fields, field)
	}

	return m
}

func (m *Message) WithContent(content string) *Message {
	m.Content = content
	return m
}

func (m *Message) Send(session *discordgo.Session, channelID string) {
	switch m.Type {
	case 1:
		m.MessageEmbed.Color = config.SuccessColor
	case -1:
		m.MessageEmbed.Color = config.ErrorColor
	default:
		m.MessageEmbed.Color = config.MessageColor
	}

	if _, err := session.ChannelMessageSendEmbed(channelID, m.MessageEmbed); err != nil {
		log.Println("error sending embed message:", err)
	}
}
