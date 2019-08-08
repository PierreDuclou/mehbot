package main

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

var usage = &discordgo.MessageEmbed{
	Title: "Liste des commandes disponibles",
	Color: 0x00dbac,
}

func init() {
	cmdnames := discordgo.MessageEmbedField{Name: "Commande", Inline: true}
	cmddesc := discordgo.MessageEmbedField{Name: "Description", Inline: true}

	for _, c := range commands {
		cmdnames.Value += fmt.Sprintf("%s%s (%s%s)\n", config.prefix, c.Name, config.prefix, c.Alias)
		cmddesc.Value += c.Description + "\n"
	}

	usage.Fields = []*discordgo.MessageEmbedField{&cmdnames, &cmddesc}
}

// Command represents a basic command provided to users
type Command struct {
	Name        string
	Alias       string
	Description string
	Usage       string
	Authroles   []string
	Session     *discordgo.Session
	MessageData *discordgo.MessageCreate
	Run         func(Command, []string) bool
}

// Execute is a façade to execute the Run() function of the command
func (c Command) Execute(args []string) {
	ok := c.Run(c, args)

	if !ok {
		err := sendEmbed(0, c.Session, c.MessageData.ChannelID, []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:  fmt.Sprintf("Utilisation de la commande **!%s**", c.Name),
				Value: c.Description + ".\n\n" + c.Usage,
			},
		})

		if err != nil {
			log.Println(err)
		}
	}
}

// Authorized determines if the given message's author is authorized to send commands to the bot
func (c Command) Authorized(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	member, err := s.GuildMember(m.GuildID, m.Author.ID)

	if err != nil {
		log.Println(err)
		return false
	}

	for _, role := range member.Roles {
		for _, authrole := range c.Authroles {
			if role == authrole {
				return true
			}
		}
	}

	return false
}

// GetCommand returns a Command instance matching the given name
func GetCommand(name string, s *discordgo.Session, m *discordgo.MessageCreate) *Command {
	for _, command := range commands {
		if name == command.Name || name == command.Alias {
			command.Session = s
			command.MessageData = m
			return &command
		}
	}
	return nil
}

var commands = []Command{
	Command{
		Name:        "help",
		Alias:       "h",
		Description: "Affiche ce message d'aide",
		Authroles: []string{
			config.baseroles["Guez"],
			config.baseroles["Guezt"],
			config.baseroles["Worms"],
		},
		Run: func(c Command, args []string) bool {
			_, err := c.Session.ChannelMessageSendEmbed(c.MessageData.ChannelID, usage)

			if err != nil {
				log.Println("error sending embed message:", err)
			}

			return true
		},
	},
	Command{
		Name:        "wplayer",
		Alias:       "wp",
		Description: "Enregistre un nouveau joueur dans la base de données",
		Usage:       "Profile : `!wp <PSEUDO> <ID DISCORD>`\n\nExemple :\n`!wp Connard 148841746661376000`",
		Authroles:   []string{config.baseroles["Superguez"]},
		Run: func(c Command, args []string) bool {
			return true
		},
	},
	Command{
		Name:  "wgame",
		Alias: "wg",
		Usage: "Format requis pour chaque ligne :\n`[*]<PSEUDO> <NOMBRE DE KILLS> <NOMBRE DE MORTS> <DÉGÂTS>`" +
			"\n\n- L'étoile désigne le vainqueur (**unique**) de la partie." +
			"\n- Les joueurs doivent avoir été enregistrés dans la base de données au préalable." +
			"\n- Les résultats de la partie doivent être écrit dans un bloc de code (entouré par trois backticks)." +
			"\n\nExemple :\n```!wg ` ` `\n*seezah 16 0 4600\nitsuped 0 8 800\ntranker 0 8 200\n` ` ` ```",
		Description: "Enregistre une nouvelle partie de worms",
		Authroles:   []string{config.baseroles["Superguez"]},
		Run: func(c Command, args []string) bool {
			return true
		},
	},
	Command{
		Name:        "wranking",
		Alias:       "wr",
		Description: "Prochainement...",
		Authroles:   []string{config.baseroles["Worms"]},
		Run: func(c Command, args []string) bool {
			return true
		},
	},
	Command{
		Name:        "wstats",
		Alias:       "ws",
		Description: "Prochainement...",
		Authroles:   []string{config.baseroles["Worms"]},
		Run: func(c Command, args []string) bool {
			return true
		},
	},
	Command{
		Name:        "youtube",
		Alias:       "yt",
		Description: "Prochainement...",
		Authroles:   nil,
		Run: func(c Command, args []string) bool {
			return true
		},
	},
}
