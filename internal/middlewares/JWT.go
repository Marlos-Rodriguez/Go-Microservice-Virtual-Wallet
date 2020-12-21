package middlewares

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/internal/environment"
	jwtware "github.com/gofiber/jwt/v2"
)

//JWTMiddleware Check the JWT of request
func JWTMiddleware() fiber.Handler {
	secrectKey, success := environment.AccessENV("SECRECT_KEY")

	if !success {
		return func(c *fiber.Ctx) error {
			return c.Status(500).JSON(&fiber.Map{
				"success": false,
				"error":   "Error in Get the Secrect key from ENV",
			})
		}
	}

	return jwtware.New(jwtware.Config{
		SigningKey: []byte(secrectKey),
	})
}
