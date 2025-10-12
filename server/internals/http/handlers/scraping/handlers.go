package scraping

import "github.com/gofiber/fiber/v2"

func ScrapingHandlers(instance *fiber.App) {
	instance.Get("/featured", FeaturedHandler)
	instance.Get("/search", SearchHandler)

	instance.Get("/:name", GetByNameHandler)
	instance.Get("/:name/:chapter", GetChapterHandler)
}
