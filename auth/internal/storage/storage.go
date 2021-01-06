package storage

import (
	"fmt"
	"log"
	"strconv"

	"github.com/jinzhu/gorm"

	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/auth/internal/environment"
	//Postgres Driver imported
	_ "github.com/lib/pq"
)

//ConnectDB connect to Postgres DB
func ConnectDB() *gorm.DB {
	var (
		host     = environment.AccessENV("DB_HOST")
		user     = environment.AccessENV("DB_USER")
		port     = environment.AccessENV("DB_PORT")
		password = environment.AccessENV("DB_PASSWORD")
		name     = environment.AccessENV("DB_NAME")
	)
	if host == "" {
		log.Fatalln("Error loading ENV")
		return nil
	}

	portInt, err := strconv.Atoi(port)

	if err != nil {
		log.Fatalln("Error in convert port to int the DB " + err.Error())
		return nil
	}

	//Connect to DB
	var DB *gorm.DB

	DB, err = gorm.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, portInt, user, password, name))

	//Check for Errors in DB
	if err != nil {
		log.Fatalf("Error in connect the DB %e", err)
		return nil
	}

	if err := DB.DB().Ping(); err != nil {
		log.Fatalln("Error in make ping the DB " + err.Error())
		return nil
	}

	if DB.Error != nil {
		log.Fatalln("Any Error in connect the DB " + err.Error())
		return nil
	}

	log.Println("DB connected")

	return DB
}
