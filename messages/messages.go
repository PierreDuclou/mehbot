package messages

import (
	"log"

	"github.com/PierreDuclou/mehbot/config"

	"github.com/bwmarrin/discordgo"
)

// MessageEmbed aims to simplify building and sending *discordgo.MessageEmbed
type MessageEmbed struct {
	MessageEmbed *discordgo.MessageEmbed
}

// Message aims to simplify sending simple messages using discordgo
type Message struct {
	Content string
}

// NewMessage returns a new *Message using the given content
func NewMessage(content string) *Message {
	return &Message{Content: content}
}

// NewMessageEmbed returns a new *MessageEmbed using the standard color
func NewMessageEmbed() *MessageEmbed {
	return &MessageEmbed{MessageEmbed: &discordgo.MessageEmbed{Color: config.StandardColor}}
}

// NewSuccessMessage returns a new *MessageEmbed using the success color
func NewSuccessMessage() *MessageEmbed {
	return &MessageEmbed{MessageEmbed: &discordgo.MessageEmbed{Color: config.SuccessColor}}
}

// NewErrorMessage returns a new *MessageEmbed using the error color
func NewErrorMessage() *MessageEmbed {
	return &MessageEmbed{MessageEmbed: &discordgo.MessageEmbed{Color: config.ErrorColor}}
}

// WithTitle binds the given title to a MessageEmbed
func (m *MessageEmbed) WithTitle(title string) *MessageEmbed {
	m.MessageEmbed.Title = title
	return m
}

// WithFields binds the given fields to a MessageEmbed
func (m *MessageEmbed) WithFields(fields ...*discordgo.MessageEmbedField) *MessageEmbed {
	for _, field := range fields {
		m.MessageEmbed.Fields = append(m.MessageEmbed.Fields, field)
	}

	return m
}

// Send sends the MessageEmbed to the given channel using the provided session
func (m *MessageEmbed) Send(session *discordgo.Session, channelID string) {
	if _, err := session.ChannelMessageSendEmbed(channelID, m.MessageEmbed); err != nil {
		log.Println("error sending embed message:", err)
	}
}

// Send sends the Message to the given channel using the provided session
func (m *Message) Send(session *discordgo.Session, channelID string) {
	if _, err := session.ChannelMessageSend(channelID, m.Content); err != nil {
		log.Println("error sending standard message:", err)
	}
}
