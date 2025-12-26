package manga

import (
	"github.com/gofiber/fiber/v2"
)

func MangaHandler(api fiber.Router) {
	api.Get("/:slug", GetMangaHandler)
	api.Get("/:slug/:chapter", GetChapterHandler)
}
