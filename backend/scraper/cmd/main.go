package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	coreutils "qaanii/scraper/internals/domain/core/core_utils"
	internal_broker "qaanii/scraper/internals/infra/broker"
	"qaanii/scraper/internals/infra/http"
	internal_utils "qaanii/scraper/internals/utils"
	"qaanii/shared/broker"
	"qaanii/shared/utils"
	"syscall"

	"github.com/gofiber/fiber/v2"
	dotenv "github.com/joho/godotenv"
)

// FIX: Replace with typed struct
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

	pool := coreutils.NewBrowserPool()
	defer pool.Cleanup()

	conn := broker.Broker(broker_url)
	defer conn.Close()

	ctx = context.WithValue(ctx, internal_utils.QUEUE_CONNECTION_KEY, conn)
	ctx = context.WithValue(ctx, internal_utils.SCRAPER_POOL_KEY, &pool)

	app := fiber.New()
	http.Router(app) // INFO: Setup HTTP debug endpoints

	internal_broker.SetupSubscribers(broker.SubscriberRequest{
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
