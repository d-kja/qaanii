package grpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"qaanii/manga/internals/constants"
	usecase "qaanii/manga/internals/domain/search/use_case"
	"qaanii/manga/internals/infra/broker"
	"qaanii/manga/internals/utils"
	base_buf "qaanii/mangabuf/gen/manga/v1"
	buf_handler "qaanii/mangabuf/gen/manga/v1/mangav1connect"
	"qaanii/shared/broker/channels"
	"qaanii/shared/broker/events"
	"time"

	"connectrpc.com/connect"
	amq "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type SearchHandler struct {
	HandlerContext *context.Context
	buf_handler.SearchServiceHandler
}

var initial_search base_buf.SearchResponse = base_buf.SearchResponse{
	Status: base_buf.RequestStatus_REQUEST_STATUS_LOADING,
	Data:   nil,
}

var err_search base_buf.SearchResponse = base_buf.SearchResponse{
	Status: base_buf.RequestStatus_REQUEST_STATUS_ERROR,
	Data:   nil,
}

const SEARCH_KEY string = "SEARCH"

var SEARCH_IDEMPOTENCY_KEY string = fmt.Sprintf("%v/IDEMPOTENCY_KEY", SEARCH_KEY)

type SearchCache struct {
	Status base_buf.RequestStatus
	Data   []*base_buf.Manga
}

func (handler SearchHandler) Search(r_ctx context.Context, request *base_buf.SearchRequest, stream *connect.ServerStream[base_buf.SearchResponse]) error {
	service := usecase.SearchByNameService{}
	ctx := *handler.HandlerContext

	redis_ch, ok := ctx.Value(constants.REDIS_URL).(*redis.Client)
	if !ok {
		return errors.New("Invalid redis client")
	}

	err := stream.Send(&initial_search)
	if err != nil {
		log.Printf("[SEARCH] - Unable to send initial status, error %+v\n", err)
		return err
	}

	idempotency_key := fmt.Sprintf("%v/%v/%v", SEARCH_IDEMPOTENCY_KEY, request.Slug, request.Id)
	idempotency, err := redis_ch.Get(r_ctx, idempotency_key).Result()
	if err == nil {
		// TODO: Abstract into a separate function
		cache := SearchCache{}

		cache_err := json.Unmarshal([]byte(idempotency), &cache)
		if cache_err == nil {
			err := stream.Send(&base_buf.SearchResponse{
				Status: cache.Status,
				Data:   cache.Data,
			})

			if err != nil {
				log.Printf("[SEARCH] - Unable to send idempotency cache, error %+v\n", err)
				return err
			}

			return nil
		}

		log.Printf("[SEARCH] - Unable to parse idempotency cache, error %+v\n", err)
	}

	// TODO: Abstract into a separate function
	query_cache, err := redis_ch.Get(r_ctx, fmt.Sprintf("%v/%v", SEARCH_KEY, request.Slug)).Result()
	if err == nil {
		cache := SearchCache{}

		cache_err := json.Unmarshal([]byte(query_cache), &cache)
		if cache_err == nil {
			err := stream.Send(&base_buf.SearchResponse{
				Status: cache.Status,
				Data:   cache.Data,
			})

			if err != nil {
				log.Printf("[SEARCH] - Unable to send saerch cache, error %+v\n", err)
				return err
			}

			return nil
		}

		log.Printf("[SEARCH] - Unable to parse search cache, error %+v\n", err)
	}

	amq_connection, ok := ctx.Value(broker.BROKER_CONNECTION).(*amq.Connection)
	if !ok {
		log.Printf("[SEARCH] - AMQ Connection not found for %v\n", events.SEARCH_MANGA_EVENT)
		return errors.New("Invalid amq connection")
	}

	channel, err := amq_connection.Channel()
	if err != nil {
		log.Printf("[SEARCH] - AMQ Channel not found for %v\n", events.SEARCH_MANGA_EVENT)
		return errors.Join(errors.New("Invalid queue publisher"), err)
	}

	queue_publisher_raw, ok := ctx.Value(events.SEARCH_MANGA_EVENT).(*events.Publisher)
	if !ok {
		log.Printf("[SEARCH] - Queue (%v) wasn't found\n", events.SEARCH_MANGA_EVENT)
		return errors.New("Invalid queue publisher")
	}

	queue_publisher := *queue_publisher_raw
	service_request := usecase.SearchByNameRequest{
		Id:     request.Id,
		Search: request.Slug,

		Channel:   channel,
		Publisher: queue_publisher,
	}

	service_response, err := service.Exec(service_request)
	if err != nil {
		stream_err := stream.Send(&err_search)
		if stream_err != nil {
			log.Printf("[SEARCH] - Unable to send service error status, error %+v\n", err)
		}

		return err
	}

	mangas := []*base_buf.Manga{}
	for _, manga := range service_response.Mangas {
		manga_buf := manga.ToProtobuf()
		mangas = append(mangas, &manga_buf)
	}

	response := base_buf.SearchResponse{
		Status: base_buf.RequestStatus_REQUEST_STATUS_COMPLETED,
		Data:   mangas,
	}

	cache_idempotency, err := json.Marshal(SearchCache{
		Status: response.Status,
		Data:   response.Data,
	})
	if err == nil {
		_, err = redis_ch.Set(r_ctx, idempotency_key, cache_idempotency, time.Duration(channels.QUEUE_TTL)*time.Second).Result()

		if err != nil {
			log.Printf("[SEARCH] - Unable to write idempotency cache, error: %v", err)
		}
	}

	err = stream.Send(&response)
	if err != nil {
		log.Printf("[SEARCH] - Unable to send service response, error %+v\n", err)
		return err
	}

	return nil
}

func SetupSearchRoute(mux *http.ServeMux, ctx *context.Context) {
	service := SearchHandler{
		HandlerContext: ctx,
	}

	path, handler := buf_handler.NewSearchServiceHandler(service)
	log.Printf("[SETUP] - Search path: %v", path)

	mux.Handle(path, utils.Middlewares(handler, ctx))
}
