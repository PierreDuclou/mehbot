package util

import (
	"log"
	"mehbot/config"

	"github.com/bwmarrin/discordgo"
)

// SendEmbed is a helper allowing to send embed messages
// messageType: -1 = error, 0 = standard, 1 = success
func SendEmbed(messageType int, session *discordgo.Session, channelID string, fields []*discordgo.MessageEmbedField) {
	message := &discordgo.MessageEmbed{}

	switch messageType {
	case 1:
		message.Color = config.SuccessColor
	case -1:
		message.Color = config.ErrorColor
	default:
		message.Color = config.MessageColor
	}

	for _, field := range fields {
		message.Fields = append(message.Fields, field)
	}

	if _, err := session.ChannelMessageSendEmbed(channelID, message); err != nil {
		log.Println("error sending embed message:", err)
	}
}
