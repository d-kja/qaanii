package grpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"qaanii/manga/internals/constants"
	usecase "qaanii/manga/internals/domain/mangas/use_case"
	"qaanii/manga/internals/infra/broker"
	"qaanii/manga/internals/utils"
	base_buf "qaanii/mangabuf/gen/manga/v1"
	buf_handler "qaanii/mangabuf/gen/manga/v1/mangav1connect"
	"qaanii/shared/broker/channels"
	"qaanii/shared/broker/events"
	"time"

	"connectrpc.com/connect"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type MangaHandler struct {
	HandlerContext *context.Context
	buf_handler.MangaServiceHandler
}

var initial_manga base_buf.GetMangaResponse = base_buf.GetMangaResponse{
	Status: base_buf.RequestStatus_REQUEST_STATUS_LOADING,
	Data:   nil,
}

var err_manga base_buf.GetMangaResponse = base_buf.GetMangaResponse{
	Status: base_buf.RequestStatus_REQUEST_STATUS_ERROR,
	Data:   nil,
}

const MANGA_KEY string = "MANGA"
var MANGA_IDEMPOTENCY_KEY string = fmt.Sprintf("%v/IDEMPOTENCY_KEY", MANGA_KEY)

type MangaCache struct {
	Status base_buf.RequestStatus
	Data   *base_buf.Manga
}

func (handler MangaHandler) GetManga(r_ctx context.Context, request *base_buf.GetMangaRequest, stream *connect.ServerStream[base_buf.GetMangaResponse]) error {
	service := usecase.GetMangaBySlugService{}
	ctx := *handler.HandlerContext

	redis_ch, ok := ctx.Value(constants.REDIS_URL).(*redis.Client)
	if !ok {
		return errors.New("Invalid redis client")
	}

	err := stream.Send(&initial_manga)
	if err != nil {
		log.Printf("[MANGA] - Unable to send initial status, error %+v\n", err)
		return err
	}

	idempotency_key := fmt.Sprintf("%v/%v/%v", MANGA_IDEMPOTENCY_KEY, request.Slug, request.Id)
	idempotency, err := redis_ch.Get(r_ctx, idempotency_key).Result()
	if err == nil {
		// TODO: Abstract into a separate function
		cache := MangaCache{}

		cache_err := json.Unmarshal([]byte(idempotency), &cache)
		if cache_err == nil {
			err := stream.Send(&base_buf.GetMangaResponse{
				Status: cache.Status,
				Data:   cache.Data,
			})

			if err != nil {
				log.Printf("[MANGA] - Unable to send idempotency cache, error %+v\n", err)
				return err
			}

			return nil
		}

		log.Printf("[SEARCH] - Unable to parse manga cache, error %+v\n", err)
	}

	amq_connection, ok := ctx.Value(broker.BROKER_CONNECTION).(*amqp091.Connection)
	if !ok {
		log.Printf("[MANGA] - AMQ Connection not found for %v\n", events.SCRAPE_MANGA_EVENT)
		return errors.New("Invalid amq connection")
	}

	channel, err := amq_connection.Channel()
	if err != nil {
		log.Printf("[MANGA] - AMQ Channel not found for %v\n", events.SEARCH_MANGA_EVENT)
		return errors.Join(errors.New("Invalid queue publisher"), err)
	}

	queue_publisher_raw, ok := ctx.Value(events.SCRAPE_MANGA_EVENT).(*events.Publisher)
	if !ok {
		log.Printf("[MANGA] - Queue (%v) wasn't found\n", events.SCRAPE_MANGA_EVENT)
		return errors.New("Invalid queue publisher")
	}

	queue_publisher := *queue_publisher_raw
	service_request := usecase.GetMangaBySlugRequest{
		Id:   request.Id,
		Slug: request.Slug,

		Channel:   channel,
		Publisher: queue_publisher,
	}

	service_response, err := service.Exec(service_request)
	if err != nil {
		stream_err := stream.Send(&err_manga)
		if stream_err != nil {
			log.Printf("[MANGA] - Unable to send service error status, error %+v\n", err)
		}

		return err
	}

	manga := service_response.Manga.ToProtobuf()
	response := base_buf.GetMangaResponse{
		Data:   &manga,
		Status: base_buf.RequestStatus_REQUEST_STATUS_COMPLETED,
	}

	cache_idempotency, err := json.Marshal(MangaCache{
		Status: response.Status,
		Data:   response.Data,
	})
	if err == nil {
		_, err = redis_ch.Set(r_ctx, idempotency_key, cache_idempotency, time.Duration(channels.QUEUE_TTL)*time.Second).Result()

		if err != nil {
			log.Printf("[MANGA] - Unable to write idempotency cache, error: %v", err)
		}
	}

	err = stream.Send(&response)
	if err != nil {
		log.Printf("[MANGA] - Unable to send service response, error %+v\n", err)
		return err
	}

	return nil
}

func SetupMangaRoute(mux *http.ServeMux, ctx *context.Context) {
	service := MangaHandler{
		HandlerContext: ctx,
	}

	path, handler := buf_handler.NewMangaServiceHandler(service)
	mux.Handle(path, utils.Middlewares(handler, ctx))
}
