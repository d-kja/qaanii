package handlers

import "github.com/gofiber/fiber/v2"

func NotFoundHandler(app *fiber.App) {
	app.Use(func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"code":    404,
			"status":  "ERROR",
			"message": "Route not found.",
		})
	})
}
