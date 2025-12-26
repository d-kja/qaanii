package search

import "github.com/gofiber/fiber/v2"

func SearchHandler(api fiber.Router) {
	api.Get("/", SearchByNameHandler)
}
