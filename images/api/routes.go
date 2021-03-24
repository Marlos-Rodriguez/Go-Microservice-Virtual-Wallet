package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/images/handlers"
	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/images/internal/middlewares"
)

func routes() *fiber.App {
	app := fiber.New()

	handler := handlers.NewImageshandlerService()

	user := app.Group("/images")

	user.Use(cors.New())

	user.Post("/:id", middlewares.JWTMiddleware(), handler.ChangeAvatar) //Recibe Avatar and changed

	return app
}
