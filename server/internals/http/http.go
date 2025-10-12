package http

import (
	"fmt"
	"os"

	"server/internals/http/handlers"
	"server/internals/http/middleware"
	"server/internals/utils"

	"github.com/gofiber/fiber/v2"
)

type HttpHandler struct{}

func (HttpHandler) Run() {
	middleware_setup := middleware.Middlewares{}
	handlers_setup := handlers.Handlers{}

	app := fiber.New(fiber.Config{
		StrictRouting: true,
		CaseSensitive: true,
	})

	middleware_setup.Consume(app)
	handlers_setup.Consume(app)

	port := os.Getenv("PORT")

	if err := app.Listen(fmt.Sprintf(":%v", port)); err != nil {
		utils.LOGGER.INFO.Fatalf("An error occurred while running application, error: %+v\n", err)
	}
}
