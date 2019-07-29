package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// App contains all the bot logic
type App struct {
	Session *discordgo.Session
}

// NewApp returns a pointer to a new App instance
func NewApp(token string) *App {
	session, err := discordgo.New("Bot " + token)

	if err != nil {
		log.Fatalln("error creating a discord session: \n\t", err)
	}

	return &App{session}
}

// Run launches the app
func (app App) Run() {
	for _, handler := range handlers {
		app.Session.AddHandler(handler)
	}

	err := app.Session.Open()

	if err != nil {
		log.Fatalln("error opening discord session: \n\t", err)
	}
}
