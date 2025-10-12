package handlers

import (
	"server/internals/http/middleware"

	"github.com/gofiber/fiber/v2"
)

func HealthHandler(instance *fiber.App) {
	instance.Get("/health", health)
}

func health(ctx *fiber.Ctx) error {
	instance := middleware.RodMiddleware{}

	scraper, page := instance.New()
	defer scraper.MustClose()

	// Test stealth
	page.MustNavigate("https://bot.sannysoft.com")
	instance.Metadata(page)

	return ctx.Status(200).SendString("Ok")
}
