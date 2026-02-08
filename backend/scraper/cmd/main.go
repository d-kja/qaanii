package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	internal_broker "qaanii/scraper/internals/infra/broker"
	"qaanii/scraper/internals/infra/http"
	"qaanii/shared/broker"
	"qaanii/shared/utils"
	"syscall"

	"github.com/gofiber/fiber/v2"
	dotenv "github.com/joho/godotenv"
)

func main() {
	signal_ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	err := dotenv.Load()
	if err != nil {
		log.Fatalf("Unable to load environment variables, error: %+v", err)
	}

	envs := utils.Envs()
	app := fiber.New()

	broker_url := envs["broker_url"]
	conn, channel := broker.Broker(broker_url)

	defer channel.Close()
	defer conn.Close()

	http.Router(app) // INFO: Setup HTTP debug endpoints

	ctx := context.Background()

	internal_broker.SetupSubscribers(broker.SubscriberRequest{
		Channel:    channel,
		Connection: conn,
		Context:    &ctx,
	})

	go func() {
		port := fmt.Sprintf(":%v", envs["port"])

		if err := app.Listen(port); err != nil {
			log.Fatalf("[SERVER] - Unable to run, error: %+v", err)
		}
	}()

	<-signal_ctx.Done()
	log.Printf("[SERVER] - Shutting down...")

	app.Shutdown()
}
