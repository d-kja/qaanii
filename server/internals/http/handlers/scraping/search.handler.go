package scraping

import (
	"server/internals/domain/scraping/services"
	"server/internals/http/middleware"
	"server/internals/utils"

	"github.com/gofiber/fiber/v2"
)

func SearchHandler(ctx *fiber.Ctx) error {
	instance := middleware.RodMiddleware{}
	query := ctx.Query("q")

	service := services.SearchMangasService{
		Scraper: instance,
	}

	response, err := service.Exec(services.SearchMangasRequest{
		Query: query,
		Ctx: ctx,
	})
	if err != nil {
		utils.LOGGER.ERROR.Printf("Unable to retrieve data, error: %+v\n", err)

		return ctx.Status(500).JSON(map[string]string{
			"status":  "ERROR",
			"message": err.Error(),
		})
	}

	return ctx.Status(200).JSON(response)
}
