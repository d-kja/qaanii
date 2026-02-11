package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	internal_broker "qaanii/scraper/internals/infra/broker"
	"qaanii/scraper/internals/infra/http"
	internal_utils "qaanii/scraper/internals/utils"
	"qaanii/shared/broker"
	"qaanii/shared/utils"
	"syscall"

	"github.com/gofiber/fiber/v2"
	dotenv "github.com/joho/godotenv"
)

var ctx = context.Background()

func main() {
	signal_ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	err := dotenv.Load()
	if err != nil {
		log.Fatalf("Unable to load environment variables, error: %+v", err)
	}

	envs := utils.Envs()
	broker_url := envs["broker_url"]

	conn, channel := broker.Broker(broker_url)
	defer channel.Close()
	defer conn.Close()

	ctx = context.WithValue(ctx, internal_utils.QUEUE_CONNECTION_KEY, conn)
	ctx = context.WithValue(ctx, internal_utils.QUEUE_CHANNEL_KEY, channel)

	app := fiber.New()
	http.Router(app) // INFO: Setup HTTP debug endpoints

	internal_broker.SetupSubscribers(broker.SubscriberRequest{
		Channel:    channel,
		Connection: conn,
		Context:    &ctx,
	})

	port := fmt.Sprintf(":%v", envs["port"])
	go handle(app, port)

	<-signal_ctx.Done()
	log.Printf("[SERVER] - Shutting down...")

	app.Shutdown()
}

func handle(app *fiber.App, port string) {
	if err := app.Listen(port); err != nil {
		log.Fatalf("[SERVER] - Unable to run, error: %+v", err)
	}
}

func reconnect() {
	// TODO: Decouple channel and connection from main, and introduce reconnection behavior required for 403 errors
}
