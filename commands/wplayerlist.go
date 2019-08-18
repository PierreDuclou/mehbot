package commands

import (
	"github.com/PierreDuclou/mehbot/config"
	"github.com/PierreDuclou/mehbot/messages"
	"github.com/PierreDuclou/mehbot/wast"

	"github.com/bwmarrin/discordgo"
)

func newWPlayerListCommand() *Command {
	cmd := Command{
		Name:        "wplayerlist",
		Alias:       "wpl",
		Description: "Affiche la liste des joueurs de worms enregistrés",
		AuthorizedRoles: []string{
			config.Roles["Guez"],
			config.Roles["Guezt"],
			config.Roles["Worms"],
		},
		Run: runWPlayerListCommand,
	}

	return &cmd
}

func runWPlayerListCommand(c Command, args []string) bool {
	players := []wast.Player{}
	wast.Db.Find(&players)
	nicknames := &discordgo.MessageEmbedField{Name: "Pseudo", Inline: true}
	ids := &discordgo.MessageEmbedField{Name: "ID Discord", Inline: true}
	messageFields := []*discordgo.MessageEmbedField{nicknames, ids}

	if len(players) == 0 {
		nicknames.Value = "-"
		ids.Value = "-"
	}

	for _, player := range players {
		nicknames.Value += player.Nickname + "\n"
		ids.Value += player.ID + "\n"
	}

	messages.NewMessageEmbed().WithFields(messageFields...).Send(c.Session, c.MessageData.ChannelID)
	return true
}
