package storage

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
	//Autoload the env
	_ "github.com/joho/godotenv/autoload"

	//Postgres Driver imported
	_ "github.com/lib/pq"
)

//ConnectDB connect to Postgres DB
func ConnectDB(service string) *gorm.DB {
	var (
		host     = os.Getenv(service + "_DB_HOST")
		user     = os.Getenv(service + "_DB_USER")
		port     = os.Getenv(service + "_DB_PORT")
		password = os.Getenv(service + "_DB_PASSWORD")
		name     = os.Getenv(service + "_DB_NAME")
	)
	if host == "" {
		log.Fatalln("Error loading ENV")
		return nil
	}

	portInt, err := strconv.Atoi(port)

	if err != nil {
		log.Fatalln("Error in connect the DB " + err.Error())
		return nil
	}

	var DB *gorm.DB

	DB, err = gorm.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, portInt, user, password, name))

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
