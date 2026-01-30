package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"qaanii/manga/internals/infra/broker"
	"qaanii/manga/internals/infra/grpc"
	"qaanii/shared/utils"

	dotenv "github.com/joho/godotenv"
)

func main() {
	err := dotenv.Load()
	if err != nil {
		log.Fatalf("Unable to load environment variables, error: %+v", err)
	}

	envs := utils.Utils{}.Envs()

	controller := grpc.GRPC{}
	ctx := context.Background()

	conn, channel := broker.Broker()
	defer channel.Close()
	defer conn.Close()

	ctx = context.WithValue(ctx, broker.BROKER_CONNECTION, conn)
	ctx = context.WithValue(ctx, broker.BROKER_CHANNEL, channel)

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

	mux, protocol := controller.Setup(&ctx)

	address := fmt.Sprintf("localhost:%v", envs["port"])
	server := http.Server{
		Addr:      address,
		Handler:   mux,
		Protocols: protocol,
	}

	log.Printf("Listening on http://%v\n", address)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Unable to run server, error: %+v", err)
	}
}
