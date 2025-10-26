package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type RecoverMiddleware struct {}
func (RecoverMiddleware) New() fiber.Handler {
	instance := recover.New()
	return instance
}
