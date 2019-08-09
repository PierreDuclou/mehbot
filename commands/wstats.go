package commands

import "mehbot/config"

func newWStatsCommand() *Command {
	cmd := Command{
		Name:        "wstats",
		Alias:       "ws",
		Description: "Prochainement...",
		AuthorizedRoles:   []string{config.Roles["Worms"]},
		Run: func(c Command, args []string) bool {
			return true
		},
	}

	return &cmd
}
