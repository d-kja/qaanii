package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	internal_broker "qaanii/manga/internals/infra/broker"
	"qaanii/manga/internals/infra/grpc"
	"qaanii/shared/broker"
	"qaanii/shared/utils"
	"syscall"
	"time"

	dotenv "github.com/joho/godotenv"
)

func main() {
	signal_ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	err := dotenv.Load()
	if err != nil {
		log.Fatalf("Unable to load environment variables, error: %+v", err)
	}

	envs := utils.Envs()

	controller := grpc.GRPC{}
	ctx := context.Background()

	broker_url := envs["broker_url"]
	conn, channel := broker.Broker(broker_url)

	defer channel.Close()
	defer conn.Close()

	ctx = context.WithValue(ctx, internal_broker.BROKER_CONNECTION, conn)
	ctx = context.WithValue(ctx, internal_broker.BROKER_CHANNEL, channel)

	internal_broker.SetupPublishers(broker.PublisherRequest{
		Channel:    channel,
		Connection: conn,
		Context:    &ctx,
	})

	internal_broker.SetupSubscribers(broker.SubscriberRequest{
		Channel:    channel,
		Connection: conn,
		Context:    &ctx,
	})

	mux, protocol := controller.Setup(&ctx)

	address := fmt.Sprintf("localhost:%v", envs["port"])
	server := http.Server{
		Addr:      address,
		Handler:   mux,
		Protocols: protocol,
	}

	go func() {
		log.Printf("[SERVER] - Listening on http://%v\n", address)

		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("[SERVER] - Unable to run, error: %+v", err)
		}
	}()

	<-signal_ctx.Done()
	log.Printf("[SERVER] - Shutting down...")

	shutdown_ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	server.Shutdown(shutdown_ctx)
}
