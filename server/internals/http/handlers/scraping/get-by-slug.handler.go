package scraping

import (
	"strings"

	"server/internals/domain/scraping/services"
	"server/internals/http/middleware"
	"server/internals/utils"

	"github.com/gofiber/fiber/v2"
)

func GetBySlugHandler(ctx *fiber.Ctx) error {
	slug := ctx.Params("slug")
	if len(strings.TrimSpace(slug)) == 0 {

		utils.LOGGER.ERROR.Printf("Slug %v is invalid\n", slug)
		return ctx.Status(400).JSON(map[string]string{
			"status":  "ERROR",
			"message": "Invalid slug was provided",
		})
	}

	scraper := middleware.RodMiddleware{}
	instance := services.GetBySlugService{
		Scraper: scraper,
	}

	payload := services.GetBySlugRequest{
		Slug: slug,
		Ctx:  ctx,
	}

	response, err := instance.Exec(payload)
	if err != nil {
		utils.LOGGER.ERROR.Printf("Get by slug handler, service error: %+v\n", err)

		return ctx.Status(500).JSON(map[string]string{
			"status":  "ERROR",
			"message": err.Error(),
		})
	}

	return ctx.Status(200).JSON(response)
}
