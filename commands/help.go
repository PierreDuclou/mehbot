package commands

import (
	"mehbot/config"
	"mehbot/util"
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
			util.SendEmbed(0, c.Session, c.MessageData.ChannelID, usage)
			return true
		},
	}

	return &cmd
}
