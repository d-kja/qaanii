package search

import (
	coreentities "qaanii/scraper/internals/domain/core/core_entities"
	usecase "qaanii/scraper/internals/domain/search/use_case"
	"qaanii/shared/utils"

	"github.com/gofiber/fiber/v2"
)

func SearchByNameHandler(ctx *fiber.Ctx) error {
	search := ctx.Query("q")
	search_len := len(search)

	if search_len == 0 || search_len >= 65 {
		return utils.Response{Status: 400, Message: "Query is required"}.GenerateResponse(ctx)
	}

	service := usecase.SearchByNameService{
		Scraper: coreentities.NewScraper(),
	}

	response, err := service.Exec(usecase.SearchByNameRequest{
		Search: search,
	})

	if err != nil {
		return utils.Response{Status: 500, Message: err.Error()}.GenerateResponse(ctx)
	}

	return utils.Response{Status: 200, Data: response.Mangas}.GenerateResponse(ctx)
}
