package main

import (
	"fmt"
	"log"
	"mehbot/wast"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var usage []*discordgo.MessageEmbedField

func init() {
	cmdnames := discordgo.MessageEmbedField{Name: "Commande", Inline: true}
	cmddesc := discordgo.MessageEmbedField{Name: "Description", Inline: true}

	for _, c := range commands {
		cmdnames.Value += fmt.Sprintf("%s%s (%s%s)\n", config.prefix, c.Name, config.prefix, c.Alias)
		cmddesc.Value += c.Description + "\n"
	}

	usage = []*discordgo.MessageEmbedField{&cmdnames, &cmddesc}
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

func (c Command) printUsage() {
	sendEmbed(0, c.Session, c.MessageData.ChannelID, []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:  fmt.Sprintf("Utilisation de la commande **!%s**", c.Name),
			Value: c.Description + ".\n\n" + c.Usage,
		},
	})
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
		Description: "Affiche la liste des commandes disponibles",
		Authroles: []string{
			config.roles["Guez"],
			config.roles["Guezt"],
			config.roles["Worms"],
		},
		Run: func(c Command, args []string) bool {
			sendEmbed(0, c.Session, c.MessageData.ChannelID, usage)
			return true
		},
	},
	Command{
		Name:        "wplayer",
		Alias:       "wp",
		Description: "Enregistre un nouveau joueur dans la base de données",
		Usage:       "Profile : `!wp <PSEUDO> <ID DISCORD>`\n\nExemple : `!wp Connard 148841746661376000`",
		Authroles:   []string{config.roles["Superguez"]},
		Run: func(c Command, args []string) bool {
			if len(args) != 2 {
				return false
			}

			nickname := args[0]
			id := args[1]

			if _, err := c.Session.GuildMember(c.MessageData.GuildID, id); err != nil {
				log.Println("error creating player:", err)
				sendEmbed(-1, c.Session, c.MessageData.ChannelID, []*discordgo.MessageEmbedField{
					&discordgo.MessageEmbedField{
						Name:  "ID Discord inconnu au bataillon",
						Value: id,
					},
				})
				return false
			}

			player := wast.NewPlayer(id, nickname)
			existing := &wast.Player{}
			db.First(&existing, &wast.Player{ID: id})

			if player.ID == existing.ID {
				sendEmbed(-1, c.Session, c.MessageData.ChannelID, []*discordgo.MessageEmbedField{
					&discordgo.MessageEmbedField{
						Name:  "Identifiant non disponible",
						Value: existing.Nickname + " " + existing.ID,
					},
				})

				log.Println("error creating new player (ID already taken):", id, nickname)
				return false
			}

			db.Create(player)
			log.Println("created new player:", player.Nickname, player.ID)
			sendEmbed(1, c.Session, c.MessageData.ChannelID, []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:  "Joueur enregistré",
					Value: nickname + " " + id,
				},
			})

			return true
		},
	},
	Command{
		Name:        "wplayerlist",
		Alias:       "wpl",
		Description: "Affiche la liste des joueurs de worms enregistrés",
		Authroles: []string{
			config.roles["Guez"],
			config.roles["Guezt"],
			config.roles["Worms"],
		},
		Run: func(c Command, args []string) bool {
			players := []wast.Player{}
			db.Find(&players)
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

			sendEmbed(0, c.Session, c.MessageData.ChannelID, messageFields)
			return true
		},
	},
	Command{
		Name:        "wgame",
		Alias:       "wg",
		Description: "Enregistre une nouvelle partie de worms",
		Usage: "Format requis pour chaque ligne :\n`[*]<PSEUDO> <NOMBRE DE VICTIMES> <NOMBRE DE MORTS> <DÉGÂTS>`" +
			"\n\n- L'étoile désigne le vainqueur (**unique**) de la partie." +
			"\n- Les joueurs doivent avoir été enregistrés dans la base de données au préalable." +
			"\n- Les résultats de la partie doivent être écrit dans un bloc de code (entouré par trois backticks)." +
			"\n\nExemple :\n```!wg ` ` `\n*seezah 16 0 4600\nitsuped 0 8 800\ntranker 0 8 200\n` ` ` ```",
		Authroles: []string{config.roles["Superguez"]},
		Run: func(c Command, args []string) bool {
			filtered := make([]string, 0)

			for _, arg := range args {
				arg = strings.Replace(arg, "```", "", -1)
				if arg != "" {
					filtered = append(filtered, arg)
				}
			}

			chunks := chunk(filtered, 4)
			var stats []*wast.Stats
			winnerPicked := false
			game := wast.Game{
				CreatedAt: time.Now(),
			}
			db.FirstOrCreate(&game, game)

			for _, chunk := range chunks {
				if len(chunk) < 4 {
					log.Println("error chunking Worms stats, chunks:", chunks)
					sendEmbed(-1, c.Session, c.MessageData.ChannelID, []*discordgo.MessageEmbedField{
						&discordgo.MessageEmbedField{
							Name:  "Erreur de format dans la commande",
							Value: fmt.Sprintf("Une ou plusieurs informations manquantes.\nDernière ligne lue : \"%s\"", strings.Join(chunk, " ")),
						},
					})
					return false
				}

				var nickname string

				winner := strings.HasPrefix(chunk[0], "*")
				if winner && winnerPicked {
					sendEmbed(-1, c.Session, c.MessageData.ChannelID, []*discordgo.MessageEmbedField{
						&discordgo.MessageEmbedField{
							Name:  "Erreur de saisie",
							Value: "Il ne peut y avoir qu'un seul vainqueur (maximum) par partie",
						},
					})
					return false
				}

				if winner {
					winnerPicked = true
					nickname = chunk[0][1:]
				} else {
					nickname = chunk[0]
				}

				player := wast.Player{}
				db.First(&player, wast.Player{Nickname: nickname})

				if player.ID == "" {
					log.Printf("player not found: \"%s\"", nickname)
					sendEmbed(-1, c.Session, c.MessageData.ChannelID, []*discordgo.MessageEmbedField{
						&discordgo.MessageEmbedField{
							Name:  "Joueur inconnu",
							Value: fmt.Sprintf("Aucun joueur nommé \"%s\" n'a été trouvé.\nUtiliser la commande **!wpl** pour afficher la liste des joueurs connus.", nickname),
						},
					})
					return false
				}

				errCallback := func(errMessage string) {
					log.Println("error parsing integer:", errMessage)
					sendEmbed(-1, c.Session, c.MessageData.ChannelID, []*discordgo.MessageEmbedField{
						&discordgo.MessageEmbedField{
							Name:  "Erreur de saisie",
							Value: errMessage,
						},
					})
				}

				kills, ok := parseInt(chunk[1], 10, 32, "Nombre de victimes non valide : "+chunk[1], errCallback)
				deaths, ok := parseInt(chunk[2], 10, 32, "Nombre de morts non valide : "+chunk[2], errCallback)
				damage, ok := parseInt(chunk[3], 10, 32, "Dégâts non valide : "+chunk[3], errCallback)

				if !ok {
					return false
				}

				stats = append(stats, &wast.Stats{
					Kills:    int(kills),
					Deaths:   int(deaths),
					Damage:   int(damage),
					Winner:   winner,
					PlayerID: player.ID,
					GameID:   game.ID,
				})
			}

			for _, stat := range stats {
				db.Create(stat)
			}

			return true
		},
	},
	Command{
		Name:        "wranking",
		Alias:       "wr",
		Description: "Prochainement...",
		Authroles:   []string{config.roles["Worms"]},
		Run: func(c Command, args []string) bool {
			return true
		},
	},
	Command{
		Name:        "wstats",
		Alias:       "ws",
		Description: "Prochainement...",
		Authroles:   []string{config.roles["Worms"]},
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
