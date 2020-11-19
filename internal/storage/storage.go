package storage

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"

	//Postgres Driver imported
	_ "github.com/lib/pq"
)

var (
	host     = os.Getenv("DB_HOST")
	port     = os.Getenv("DB_PORT")
	user     = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASSWORD")
	name     = os.Getenv("DB_NAME")
)

//ConnectDB connect to Postgres DB
func ConnectDB() *gorm.DB {
	portInt, err := strconv.Atoi(port)

	if err != nil {
		log.Fatalln("Error in connect the DB " + err.Error())
		return nil
	}
	DB, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, portInt, user, password, name))

	if err != nil {
		log.Fatalln("Error in connect the DB " + err.Error())
		return nil
	}

	if err := DB.DB().Ping(); err != nil {
		log.Fatalln("Error in connect the DB " + err.Error())
		return nil
	}

	if DB.Error != nil {
		log.Fatalln("Error in connect the DB " + err.Error())
		return nil
	}

	log.Println("DB connected")

	return DB
}
