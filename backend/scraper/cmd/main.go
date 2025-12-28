package main

import (
	"fmt"
	"log"
	"qaanii/scraper/internals/infra/broker"
	"qaanii/scraper/internals/infra/http"
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

	// INFO: Setup HTTP debug endpoints

	conn, channel := broker.Broker(app)
	defer channel.Close()
	defer conn.Close()

	http.Router(app)

	broker.SetupConsumers(broker.ConsumerRequest{
		Channel:    channel,
		Connection: conn,
	})

	port := fmt.Sprintf(":%v", envs["port"])
	if err := app.Listen(port); err != nil {
		log.Fatalf("Unable to run server, error: %+v", err)
	}
}
