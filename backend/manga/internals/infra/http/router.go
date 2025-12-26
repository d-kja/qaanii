package http

import (
	"qaanii/manga/internals/infra/http/manga"
	"qaanii/manga/internals/infra/http/search"

	"github.com/gofiber/fiber/v2"
)

const PREFIX string = "/api/v1"

func Router(app *fiber.App) {
	app.Route(PREFIX, func(api fiber.Router) {
		api.Route("/search", search.SearchHandler)
		api.Route("/manga", manga.MangaHandler)
	})
}
