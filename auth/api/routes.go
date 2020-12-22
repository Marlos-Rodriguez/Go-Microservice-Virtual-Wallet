package api

import (
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jinzhu/gorm"

	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/auth/handlers"
	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/internal/middlewares"
)

func routes(DB *gorm.DB, RDB *redis.Client) *fiber.App {
	app := fiber.New()

	handler := handlers.NewAuthHandlerService(DB, RDB)

	auth := app.Group("/auth")

	auth.Use(cors.New())

	auth.Post("/register", handler.Register)
	auth.Post("/login", handler.Login)
	auth.Put("/reactivate", handler.ReactivateUser)
	auth.Delete("/delete", middlewares.JWTMiddleware(), handler.DeactivateUser)

	return app
}
