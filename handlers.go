package main

import (
	"log"
	"mehbot/commands"
	"mehbot/config"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var handlers = []interface{}{
	messageCreate,
}

// messageCreate is called whenever a message is pushed in a known channel
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !strings.HasPrefix(m.Message.Content, config.Prefix) {
		return
	}

	cmdname := strings.Fields(m.Message.Content)[0][1:]
	args := strings.Fields(m.Message.Content)[1:]
	command := commands.GetCommand(cmdname, s, m)

	if command == nil {
		log.Printf("unknown command \"%s\" called by user %s", cmdname, m.Author.String())
		return
	}

	if !command.Authorized(s, m) {
		log.Printf("unauthorized user \"%s\" typed:\"%s\"\n", m.Author.String(), m.Content)
		return
	}

	log.Printf("command called by user \"%s\": %s\n", m.Author.String(), m.Content)
	command.Execute(args)
}
