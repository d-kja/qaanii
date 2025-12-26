package manga

import (
	coreentities "qaanii/scraper/internals/domain/core/core_entities"
	usecase "qaanii/scraper/internals/domain/mangas/use_case"
	"qaanii/shared/utils"

	"github.com/gofiber/fiber/v2"
)

func GetMangaHandler(ctx *fiber.Ctx) error {
	slug := ctx.Params("slug")
	slug_len := len(slug)

	if slug_len == 0 {
		return utils.Response{Status: 400, Message: "Slug is required"}.GenerateResponse(ctx)
	}

	service := usecase.GetMangaBySlugService{
		Scraper: coreentities.NewScraper(),
	}

	response, err := service.Exec(usecase.GetMangaBySlugRequest{
		Slug: slug,
	})

	if err != nil {
		return utils.Response{Status: 500, Message: err.Error()}.GenerateResponse(ctx)
	}

	return utils.Response{Status: 200, Data: response.Manga}.GenerateResponse(ctx)
}
