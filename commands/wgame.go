package commands

import (
	"fmt"
	"log"
	"mehbot/config"
	"mehbot/util"
	"mehbot/wast"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func newWGameCommand() *Command {
	cmd := Command{
		Name:        "wgame",
		Alias:       "wg",
		Description: "Enregistre une nouvelle partie de worms",
		Usage: "Format requis pour chaque ligne :\n`[*]<PSEUDO> <NOMBRE DE VICTIMES> <NOMBRE DE MORTS> <DÉGÂTS>`" +
			"\n\n- L'étoile désigne le vainqueur (**unique**) de la partie." +
			"\n- Les joueurs doivent avoir été enregistrés dans la base de données au préalable." +
			"\n- Les résultats de la partie doivent être écrit dans un bloc de code (entouré par trois backticks)." +
			"\n\nExemple :\n```!wg ` ` `\n*seezah 16 0 4600\nitsuped 0 8 800\ntranker 0 8 200\n` ` ` ```",
		Authroles: []string{config.Roles["Superguez"]},
		Run:       runWGameCommand,
	}

	return &cmd
}

func runWGameCommand(c Command, args []string) bool {
	filtered := make([]string, 0)

	for _, arg := range args {
		arg = strings.Replace(arg, "```", "", -1)
		if arg != "" {
			filtered = append(filtered, arg)
		}
	}

	chunks := util.Chunk(filtered, 4)
	var stats []*wast.Stats
	winnerPicked := false
	game := wast.Game{
		CreatedAt: time.Now(),
	}
	wast.Db.FirstOrCreate(&game, game)

	for _, chunk := range chunks {
		if len(chunk) < 4 {
			log.Println("error chunking Worms stats, chunks:", chunks)
			util.SendEmbed(-1, c.Session, c.MessageData.ChannelID, []*discordgo.MessageEmbedField{
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
			util.SendEmbed(-1, c.Session, c.MessageData.ChannelID, []*discordgo.MessageEmbedField{
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
		wast.Db.First(&player, wast.Player{Nickname: nickname})

		if player.ID == "" {
			log.Printf("player not found: \"%s\"", nickname)
			util.SendEmbed(-1, c.Session, c.MessageData.ChannelID, []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:  "Joueur inconnu",
					Value: fmt.Sprintf("Aucun joueur nommé \"%s\" n'a été trouvé.\nUtiliser la commande **!wpl** pour afficher la liste des joueurs connus.", nickname),
				},
			})
			return false
		}

		errCallback := func(errMessage string) {
			log.Println("error parsing integer:", errMessage)
			util.SendEmbed(-1, c.Session, c.MessageData.ChannelID, []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:  "Erreur de saisie",
					Value: errMessage,
				},
			})
		}

		kills, ok := util.ParseInt(chunk[1], 10, 32, "Nombre de victimes non valide : "+chunk[1], errCallback)
		deaths, ok := util.ParseInt(chunk[2], 10, 32, "Nombre de morts non valide : "+chunk[2], errCallback)
		damage, ok := util.ParseInt(chunk[3], 10, 32, "Dégâts non valide : "+chunk[3], errCallback)

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
		wast.Db.Create(stat)
	}

	return true
}
