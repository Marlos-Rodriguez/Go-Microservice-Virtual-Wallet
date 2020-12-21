package api

import (
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jinzhu/gorm"

	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/internal/middlewares"
	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/user/handlers"
)

func routes(DB *gorm.DB, RDB *redis.Client) *fiber.App {
	app := fiber.New()

	handler := handlers.NewUserhandlerService(DB, RDB)

	user := app.Group("/user")

	user.Use(cors.New())

	user.Get("/:id", middlewares.JWTMiddleware(), handler.GetUser)                           //Get Basic user info
	user.Get("/all/:id", middlewares.JWTMiddleware(), handler.GetProfileUser)                //Get Profile User Info
	user.Get("/relation/:id/:page", middlewares.JWTMiddleware(), handler.GetRelations)       //Get relations of user
	user.Put("/:id", middlewares.JWTMiddleware(), handler.ModifyUser)                        //Modify the user info
	user.Post("/add", middlewares.JWTMiddleware(), handler.CreateRelation)                   //Create a new relation
	user.Delete("/relation/delete", middlewares.JWTMiddleware(), handler.DeactivateRelation) //Delete a relation

	return app
}
