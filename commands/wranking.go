package commands

import "mehbot/config"

func newWRankingCommand() *Command {
	cmd := Command{
		Name:        "wranking",
		Alias:       "wr",
		Description: "Prochainement...",
		AuthorizedRoles:   []string{config.Roles["Worms"]},
		Run: func(c Command, args []string) bool {
			return true
		},
	}

	return &cmd
}