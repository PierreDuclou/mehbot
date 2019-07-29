package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app := NewApp(os.Getenv("MEHBOT_TOKEN"))
	app.Run()
	defer app.Session.Close()
	fmt.Println("Mehbot iz running! Press CTRL-C to exit...")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
