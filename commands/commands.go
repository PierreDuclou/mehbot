package commands

import (
	"fmt"
	"log"
	"mehbot/config"
	"mehbot/util"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Commands is the slice containing all the available commands
var Commands []*Command
var usage []*discordgo.MessageEmbedField

func init() {
	addCommands(
		newHelpCommand(),
		newWGameCommand(),
		newWPlayerCommand(),
		newWPlayerListCommand(),
		newWRankingCommand(),
		newWStatsCommand(),
	)
	cmdnames := discordgo.MessageEmbedField{Name: "Commande", Inline: true}
	cmddesc := discordgo.MessageEmbedField{Name: "Description", Inline: true}

	for _, c := range Commands {
		cmdnames.Value += fmt.Sprintf("%s%s (%s%s)\n", config.Prefix, c.Name, config.Prefix, c.Alias)
		cmddesc.Value += c.Description + "\n"
	}

	usage = []*discordgo.MessageEmbedField{&cmdnames, &cmddesc}
}

func addCommands(commands ...*Command) {
	for _, command := range commands {
		Commands = append(Commands, command)
	}
}

// GetCommand returns a Command instance matching the given name
func GetCommand(name string, s *discordgo.Session, m *discordgo.MessageCreate) *Command {
	for _, command := range Commands {
		if name == command.Name || name == command.Alias {
			command.Session = s
			command.MessageData = m
			return command
		}
	}

	return nil
}

// Command represents a basic command provided to users
type Command struct {
	Name            string
	Alias           string
	Description     string
	Usage           string
	AuthorizedRoles []string
	Session         *discordgo.Session
	MessageData     *discordgo.MessageCreate
	Run             func(Command, []string) bool
}

func (c Command) printUsage() {
	util.NewMessage(0).WithFields(&discordgo.MessageEmbedField{
		Name:  fmt.Sprintf("Utilisation de la commande **!%s**", c.Name),
		Value: c.Description + ".\n\n" + c.Usage,
	}).Send(c.Session, c.MessageData.ChannelID)
}

// Execute is a façade to execute the Run() function of the command
func (c Command) Execute(args []string) {
	if len(args) > 0 {
		if args[0] == strings.ToLower("help") {
			c.printUsage()
			return
		}
	}

	ok := c.Run(c, args)

	if !ok {
		log.Printf("command %s failed to execute, args: %s\n", c.Name, args)
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
		for _, authrole := range c.AuthorizedRoles {
			if role == authrole {
				return true
			}
		}
	}

	return false
}
