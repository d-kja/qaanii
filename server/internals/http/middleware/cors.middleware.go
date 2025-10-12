package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type CorsMiddleware struct {}
func (CorsMiddleware) New() fiber.Handler {
	instance := cors.New()
	return instance
}
