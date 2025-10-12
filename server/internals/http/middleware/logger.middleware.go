package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type LoggerMiddleware struct{}
func (LoggerMiddleware) New() fiber.Handler {
	instance := logger.New()
	return instance
}
