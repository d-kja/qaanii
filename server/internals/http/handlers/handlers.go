package handlers

import (
	"server/internals/http/handlers/scraping"

	"github.com/gofiber/fiber/v2"
)

type Handlers struct {}
func (Handlers) Consume(instance *fiber.App) {
	HealthHandler(instance)

	scraping.ScrapingHandlers(instance)	

	// Fallback
	NotFoundHandler(instance)
}
