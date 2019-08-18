package commands

import (
	"github.com/PierreDuclou/mehbot/config"
	"github.com/PierreDuclou/mehbot/messages"
)

func newHelpCommand() *Command {
	cmd := Command{
		Name:        "help",
		Alias:       "h",
		Description: "Affiche la liste des commandes disponibles",
		AuthorizedRoles: []string{
			config.Roles["Guez"],
			config.Roles["Guezt"],
			config.Roles["Worms"],
		},
		Run: func(c Command, args []string) bool {
			messages.NewMessageEmbed().WithTitle("Liste des commandes").WithFields(usage...).Send(c.Session, c.MessageData.ChannelID)
			return true
		},
	}

	return &cmd
}
