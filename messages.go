package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// messageType: -1 = error, 0 = standard, 1 = success
func sendEmbed(messageType int, session *discordgo.Session, channelID string, fields []*discordgo.MessageEmbedField) {
	message := &discordgo.MessageEmbed{}

	switch messageType {
	case 1:
		message.Color = config.successColor
	case -1:
		message.Color = config.errorColor
	default:
		message.Color = config.messageColor
	}

	for _, field := range fields {
		message.Fields = append(message.Fields, field)
	}

	if _, err := session.ChannelMessageSendEmbed(channelID, message); err != nil {
		log.Println("error sending embed message:", err)
	}
}
