package storage

import (
	"log"

	"github.com/jinzhu/gorm"

	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/internal/environment"
	//Postgres Driver imported
	_ "github.com/lib/pq"
)

//ConnectDB connect to Postgres DB
func ConnectDB(service string) *gorm.DB {
	//Variables for DB
	dbAddr, success := environment.AccessENV(service + "_DB")
	if !success {
		log.Fatalln("Error loading ENV")
		return nil
	}

	//Connect to DB
	var DB *gorm.DB

	DB, err := gorm.Open(dbAddr)

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
