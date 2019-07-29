package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// Command represents a basic command provided to users
type Command struct {
	Name        string
	Alias       string
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

// GetCommand returns the corresponding Command depending on the given name
func GetCommand(name string, s *discordgo.Session, m *discordgo.MessageCreate) *Command {
	for _, command := range commands {
		if name == command.Name || command.Alias == name {
			command.Session = s
			command.MessageData = m
			return &command
		}
	}
	return nil
}

var commands = []Command{}
