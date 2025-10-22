package scraping

import "github.com/gofiber/fiber/v2"

func ScrapingHandlers(instance *fiber.App) {
	instance.Get("/featured", FeaturedHandler) // TODO:
	instance.Get("/search", SearchHandler)

	instance.Get("/manga/:slug", GetBySlugHandler)
	instance.Get("/manga/:slug/:chapter", GetChapterHandler) // TODO:
}
