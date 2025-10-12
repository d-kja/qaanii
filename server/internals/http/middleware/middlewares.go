package middleware

import "github.com/gofiber/fiber/v2"

type Middlewares struct{}
func (Middlewares) Consume(instance *fiber.App) {
	instance.Use(
		LoggerMiddleware{}.New(),
		CorsMiddleware{}.New(),
	)
}
