package search

import (
	usecase "qaanii/manga/internals/domain/search/use_case"
	"qaanii/shared/utils"

	"github.com/gofiber/fiber/v2"
)

const LIMIT_SEARCH_LEN int = 65

func SearchByNameHandler(ctx *fiber.Ctx) error {
	search := ctx.Query("q")

	search_len := len(search)
	if search_len == 0 || search_len >= LIMIT_SEARCH_LEN {
		return utils.Response{Status: 400, Message: "Query is required"}.GenerateResponse(ctx)
	}

	service := usecase.SearchByNameService{}
	request := usecase.SearchByNameRequest{
		Search: search,
	}

	response, err := service.Exec(request)
	if err != nil {
		return utils.Response{Status: 500, Message: err.Error()}.GenerateResponse(ctx)
	}

	return utils.Response{Status: 200, Data: response.Mangas}.GenerateResponse(ctx)
}
