package utils

import (
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`

	Data any `json:"data"`
}

func (self Response) GenerateResponse(ctx *fiber.Ctx) error {
	is_successful := self.Status >= 200 && self.Status <= 299
	if len(self.Message) == 0 && is_successful {
		self.Message = "Request finished successfully!" // Just to avoid empty messages
	}

	return ctx.Status(self.Status).JSON(self)
}
