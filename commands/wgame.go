package commands

import (
	"fmt"
	"log"
	"mehbot/config"
	"mehbot/messages"
	"mehbot/util"
	"mehbot/wast"
	"strconv"
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
			"\n- Les statistiques de la partie peuvent être écrit dans un bloc de code (entouré par trois backticks) pour simplifier leur saisie." +
			"\n\nExemple :\n```!wg ` ` `\n*seezah 16 0 4600\nitsuped 0 8 800\ntranker 0 8 200\n` ` ` ```",
		AuthorizedRoles: []string{config.Roles["Superguez"]},
		Run:             runWGameCommand,
	}

	return &cmd
}

func runWGameCommand(c Command, args []string) bool {
	filtered := filterArgs(args)
	chunks := util.Chunk(filtered, 4)
	var stats []*wast.Stats
	winnerPicked := false
	game := wast.Game{
		CreatedAt: time.Now(),
	}
	wast.Db.FirstOrCreate(&game, game)

	for _, chunk := range chunks {
		if len(chunk) < 4 {
			throwChunkingError(c, chunk)
			return false
		}

		var nickname string

		winner := strings.HasPrefix(chunk[0], "*")
		if winner && winnerPicked {
			throwWinnerError(c)
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
			throwPlayerNotFoundError(c, nickname)
			return false
		}

		kills, deaths, damage, ok := parseScore(&[3]string{chunk[1], chunk[2], chunk[3]})

		if !ok {
			throwParseIntError(c, chunk)
			return false
		}

		stat := wast.NewStats(int(kills), int(deaths), int(damage), winner, player.ID, game.ID)
		stat.Player = player
		stats = append(stats, stat)
	}

	for _, stat := range stats {
		wast.Db.Create(stat)
	}

	notifyCommandSuccess(c, stats)
	return true
}

func filterArgs(args []string) []string {
	filtered := make([]string, 0)

	for _, arg := range args {
		arg = strings.Replace(arg, "```", "", -1)
		if arg != "" {
			filtered = append(filtered, arg)
		}
	}

	return filtered
}

func throwChunkingError(c Command, chunk []string) {
	log.Println("error chunking Worms stats, last chunk:", chunk)
	messages.NewSuccessMessage().WithFields(&discordgo.MessageEmbedField{
		Name:  "Erreur de format dans la commande",
		Value: fmt.Sprintf("Une ou plusieurs informations manquantes dans la ligne :\n\"%s\"", strings.Join(chunk, " ")),
	}).Send(c.Session, c.MessageData.ChannelID)
}

func throwWinnerError(c Command) {
	messages.NewErrorMessage().WithFields(&discordgo.MessageEmbedField{
		Name:  "Erreur de saisie",
		Value: "Il ne peut y avoir qu'un seul vainqueur (maximum) par partie",
	}).Send(c.Session, c.MessageData.ChannelID)
}

func throwPlayerNotFoundError(c Command, nickname string) {
	log.Printf("player not found: \"%s\"", nickname)
	messages.NewErrorMessage().WithFields(&discordgo.MessageEmbedField{
		Name:  "Joueur inconnu",
		Value: fmt.Sprintf("Aucun joueur nommé \"%s\" n'a été trouvé.\nUtiliser la commande **!wpl** pour afficher la liste des joueurs connus.", nickname),
	}).Send(c.Session, c.MessageData.ChannelID)
}

func throwParseIntError(c Command, chunk []string) {
	log.Printf("error parsing integer in chunk: \"%s\"\n", chunk)
	messages.NewErrorMessage().WithFields(&discordgo.MessageEmbedField{
		Name:  "Erreur de saisie",
		Value: fmt.Sprintf("Une ou plusieurs valeurs numériques non valides sur la ligne :\n\"%s\"", strings.Join(chunk, " ")),
	}).Send(c.Session, c.MessageData.ChannelID)
}

func parseScore(buf *[3]string) (int, int, int, bool) {
	var values [3]int

	for i, value := range *buf {
		intval, err := strconv.ParseInt(value, 10, 32)

		if err != nil {
			return 0, 0, 0, false
		}

		values[i] = int(intval)
	}

	return values[0], values[1], values[2], true
}

func notifyCommandSuccess(c Command, stats []*wast.Stats) {
	log.Printf("new game inserted: %v\n", stats)
	format := "%-12s %-12v %-12v %-12v\n"
	content := ">>> **Nouvelle partie enregistrée :**\n"
	content += fmt.Sprintf("```"+format, "Joueur", "Victimes", "Morts", "Dégâts")
	content += strings.Repeat("-", 45) + "\n"

	for _, stat := range stats {
		victory := ""

		if stat.Winner {
			victory = "*"
		}

		content += fmt.Sprintf(format, victory+stat.Player.Nickname, stat.Kills, stat.Deaths, stat.Damage)
	}

	content += "```"
	messages.NewMessage(content).Send(c.Session, c.MessageData.ChannelID)
}
