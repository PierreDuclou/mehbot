package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

func init() {
	var err error
	host := os.Getenv("PSQL_HOST")
	port := os.Getenv("PSQL_PORT")
	user := os.Getenv("PSQL_USER")
	dbname := os.Getenv("PSQL_DBNAME")
	password := os.Getenv("PSQL_PASSWD")
	args := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, port, user, dbname, password)
	db, err = gorm.Open("postgres", args)

	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	app := NewApp(os.Getenv("MEHBOT_TOKEN"))
	app.Run()
	defer app.Session.Close()
	fmt.Println("Mehbot iz running! Press CTRL-C to exit...")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
