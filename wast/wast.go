package wast

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // providing the postgres driver
)

// Db is the connection to the wast database
var Db *gorm.DB

func init() {
	var err error
	host := os.Getenv("PSQL_HOST")
	port := os.Getenv("PSQL_PORT")
	user := os.Getenv("PSQL_USER")
	dbname := os.Getenv("PSQL_DBNAME")
	password := os.Getenv("PSQL_PASSWD")
	args := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, port, user, dbname, password)
	Db, err = gorm.Open("postgres", args)

	if err != nil {
		log.Fatalln(err)
	}
}
