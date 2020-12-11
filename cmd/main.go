package main

import (
	"log"
	//Autoload the env
	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/internal/storage"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	DB := storage.ConnectDB("USER")

	log.Println(DB.Rows())
}
