package main

import (
	"context"
	"fmt"
	"log"
	"qaanii/manga/internals/infra/broker"
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

	ctx := context.Background()
	app := fiber.New()

	conn, channel := broker.Broker(app)
	defer channel.Close()
	defer conn.Close()

	http.Router(app)

	broker.SetupPublishers(broker.PublisherRequest{
		Channel:    channel,
		Connection: conn,
		Context:    &ctx,
	})

	broker.SetupSubscribers(broker.SubscriberRequest{
		Channel:    channel,
		Connection: conn,
		Context:    &ctx,
	})

	port := fmt.Sprintf(":%v", envs["port"])
	if err := app.Listen(port); err != nil {
		log.Fatalf("Unable to run server, error: %+v", err)
	}
}
