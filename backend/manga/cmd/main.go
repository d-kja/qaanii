package main

import (
	"database/sql"
	_ "modernc.org/sqlite" // Relies on side effects with init function

	"context"
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"qaanii/manga/internals/constants"
	internal_broker "qaanii/manga/internals/infra/broker"
	"qaanii/manga/internals/infra/grpc"
	"qaanii/shared/broker"
	"qaanii/shared/utils"
	internal_utils "qaanii/manga/internals/utils"
	"syscall"
	"time"

	dotenv "github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

//go:embed schemas.sql
var schema string
var ctx = context.Background()

const buffer_size int = 1024 * 1024 * 15 // 15 MiB

func main() {
	signal_ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	
	defer internal_utils.Recover()

	err := dotenv.Load()
	if err != nil {
		log.Fatalf("Unable to load environment variables, error: %+v", err)
	}

	envs := utils.Envs()

	controller := grpc.GRPC{}

	redis_url := envs["redis_url"]
	broker_url := envs["broker_url"]
	database_url := envs["database_url"]

	db, err := sql.Open("sqlite", database_url)
	if err != nil {
		log.Fatalf("Unable to connect to database URL, error: %+v", err)
	}

	if _, err := db.ExecContext(ctx, schema); err != nil {
		log.Fatalf("Unable to execute database schema query, error: %+v", err)
	}

	opts, err := redis.ParseURL(redis_url)
	if err != nil {
		log.Fatalf("Unable to parse redis URL, error: %+v", err)
	}

	opts.ReadBufferSize = buffer_size
	opts.WriteBufferSize = buffer_size

	conn := broker.Broker(broker_url)
	defer conn.Close()

	redis_ch := redis.NewClient(opts)
	defer redis_ch.Close()

	// Cry all u want, I don't want to use a struct or dep injection framework, f u bish
	ctx = context.WithValue(ctx, constants.DATABASE_URL, db)
	ctx = context.WithValue(ctx, constants.REDIS_URL, redis_ch)
	ctx = context.WithValue(ctx, internal_broker.BROKER_CONNECTION, conn)

	internal_broker.SetupPublishers(broker.PublisherRequest{
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

	go handle(&server, address)
	<-signal_ctx.Done()

	log.Printf("[SERVER] - Shutting down...")
	shutdown_ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	server.Shutdown(shutdown_ctx)
}

func handle(server *http.Server, address string) {
	log.Printf("[SERVER] - Listening on http://%v\n", address)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("[SERVER] - Unable to run, error: %+v", err)
	}
}
