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
	"qaanii/shared/broker/events"

	"connectrpc.com/connect"
	amq "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type SearchHandler struct {
	HandlerContext *context.Context
	buf_handler.SearchServiceHandler
}

var initial_status base_buf.SearchResponse = base_buf.SearchResponse{
	Status: base_buf.RequestStatus_REQUEST_STATUS_LOADING,
	Data:   nil,
}

var err_response base_buf.SearchResponse = base_buf.SearchResponse{
	Status: base_buf.RequestStatus_REQUEST_STATUS_ERROR,
	Data:   nil,
}

const SEARCH_KEY string = "SEARCH"

var IDEMPOTENCY_KEY string = fmt.Sprintf("%v/IDEMPOTENCY_KEY", SEARCH_KEY)

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

	err := stream.Send(&initial_status)
	if err != nil {
		log.Printf("[SEARCH] - Unable to send initial status, error %+v\n", err)
		return err
	}

	// TODO: Abstract into a separate function
	idempotency, err := redis_ch.Get(r_ctx, fmt.Sprintf("%v/%v/%v", IDEMPOTENCY_KEY, request.Slug, request.Id)).Result()
	if err == nil {
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

	channel, ok := ctx.Value(broker.BROKER_CHANNEL).(*amq.Channel)
	if !ok {
		log.Printf("[SEARCH] - AMQ Channel not found for %v\n", events.SEARCH_MANGA_EVENT)
		return errors.New("Invalid queue publisher")
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
		err := stream.Send(&err_response)

		if err != nil {
			log.Printf("[SEARCH] - Unable to send service error status, error %+v\n", err)
			return err
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
	mux.Handle(path, utils.Middlewares(handler, ctx))
}
