package api

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func createServer(app *fiber.App) {
	//Get the Port from ENV
	PORT := os.Getenv("USER_PORT")

	if PORT == "" {
		PORT = "3000"
	}

	log.Println("Server running in Port: " + PORT)

	log.Fatal(app.Listen(":" + PORT))
}
