package commands

import (
	"log"
	"mehbot/config"
	"mehbot/util"
	"mehbot/wast"

	"github.com/bwmarrin/discordgo"
)

func newWPlayerCommand() *Command {
	cmd := Command{
		Name:        "wplayer",
		Alias:       "wp",
		Description: "Enregistre un nouveau joueur dans la base de données",
		Usage:       "Profile : `!wp <PSEUDO> <ID DISCORD>`\n\nExemple : `!wp Connard 148841746661376000`",
		AuthorizedRoles:   []string{config.Roles["Superguez"]},
		Run:         runWPlayerCommand,
	}

	return &cmd
}

func runWPlayerCommand(c Command, args []string) bool {
	if len(args) != 2 {
		return false
	}

	nickname := args[0]
	id := args[1]

	if _, err := c.Session.GuildMember(c.MessageData.GuildID, id); err != nil {
		log.Println("error creating player:", err)
		util.SendEmbed(-1, c.Session, c.MessageData.ChannelID, []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:  "ID Discord inconnu au bataillon",
				Value: id,
			},
		})
		return false
	}

	player := wast.NewPlayer(id, nickname)
	existing := &wast.Player{}
	wast.Db.First(&existing, &wast.Player{ID: id})

	if player.ID == existing.ID {
		util.SendEmbed(-1, c.Session, c.MessageData.ChannelID, []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:  "Identifiant non disponible",
				Value: existing.Nickname + " " + existing.ID,
			},
		})

		log.Println("error creating new player (ID already taken):", id, nickname)
		return false
	}

	wast.Db.Create(player)
	log.Println("created new player:", player.Nickname, player.ID)
	util.SendEmbed(1, c.Session, c.MessageData.ChannelID, []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:  "Joueur enregistré",
			Value: nickname + " " + id,
		},
	})

	return true
}
