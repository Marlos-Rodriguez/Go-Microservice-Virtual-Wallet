package api

import (
	"log"

	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/internal/environment"
	"github.com/gofiber/fiber/v2"
)

func createServer(app *fiber.App) {
	//Get the Port from ENV
	PORT := environment.AccessENV("TRANSACTION_PORT")

	if PORT == "" {
		PORT = "3002"
	}

	log.Println("Server running in Port: " + PORT)

	log.Fatal(app.Listen(":" + PORT))
}
