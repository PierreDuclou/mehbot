package main

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

var usage = &discordgo.MessageEmbed{
	Title: "Liste des commandes disponibles",
	Color: 0x4bfcc4,
}

func init() {
	cmdnames := discordgo.MessageEmbedField{Name: "Commande", Inline: true}
	cmddesc := discordgo.MessageEmbedField{Name: "Description", Inline: true}

	for _, command := range commands {
		cmdnames.Value += fmt.Sprintf("%s%s (%s%s)\n", config.prefix, command.Name, config.prefix, command.Alias)
		cmddesc.Value += command.Description + "\n"
	}

	usage.Fields = []*discordgo.MessageEmbedField{&cmdnames, &cmddesc}
}

// Command represents a basic command provided to users
type Command struct {
	Name        string
	Alias       string
	Description string
	Authroles   []string
	Session     *discordgo.Session
	MessageData *discordgo.MessageCreate
	Run         func(Command, []string)
}

// Execute is a façade to execute the Run() function of the command
func (c Command) Execute(args []string) {
	c.Run(c, args)
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
		Authroles:   config.baseroles,
		Run: func(c Command, args []string) {
			_, err := c.Session.ChannelMessageSendEmbed(c.MessageData.ChannelID, usage)

			if err != nil {
				log.Println("error sending message embed:", err)
			}
		},
	},
	Command{
		Name:        "worms",
		Alias:       "w",
		Description: "Prochainement...",
		Authroles:   config.baseroles,
		Run: func(c Command, args []string) {

		},
	},
	Command{
		Name:        "youtube",
		Alias:       "yt",
		Description: "Prochainement...",
		Authroles:   config.baseroles,
		Run: func(c Command, args []string) {

		},
	},
}
