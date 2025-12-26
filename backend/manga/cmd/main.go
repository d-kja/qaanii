package main

import (
	"fmt"
	"log"
	"qaanii/manga/internals/infra/http"
	"qaanii/shared/utils"

	"github.com/gofiber/fiber/v2"
	dotenv "github.com/joho/godotenv"
)

func main() {
	err := dotenv.Load()
	if err != nil {
		log.Fatalf("Unable to load environment variables, error: %+v", err)
	}

	envs := utils.Utils{}.Envs()

	app := fiber.New()
	http.Router(app)

	port := fmt.Sprintf(":%v", envs["port"])
	if err := app.Listen(port); err != nil {
		log.Fatalf("Unable to run server, error: %+v", err)
	}
}
