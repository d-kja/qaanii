package manga

import (
	coreentities "qaanii/scraper/internals/domain/core/core_entities"
	usecase "qaanii/scraper/internals/domain/mangas/use_case"
	"qaanii/shared/utils"

	"github.com/gofiber/fiber/v2"
)

func GetChapterHandler(ctx *fiber.Ctx) error {
	slug := ctx.Params("slug")
	chapter := ctx.Params("chapter")

	slug_len := len(slug)
	chapter_len := len(chapter)

	if slug_len == 0 || chapter_len == 0 {
		return utils.Response{Status: 400, Message: "Params are required"}.GenerateResponse(ctx)
	}

	service := usecase.GetMangaChapterService{
		Scraper: coreentities.NewScraper(),
	}

	response, err := service.Exec(usecase.GetMangaChapterRequest{
		Slug:    slug,
		Chapter: chapter,
	})

	if err != nil {
		return utils.Response{Status: 500, Message: err.Error()}.GenerateResponse(ctx)
	}

	return utils.Response{Status: 200, Data: response.Chapter}.GenerateResponse(ctx)
}
