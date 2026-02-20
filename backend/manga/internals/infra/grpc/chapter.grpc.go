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

type ChapterHandler struct {
	HandlerContext *context.Context
	buf_handler.ChapterServiceHandler
}

var initial_chapter base_buf.GetChapterResponse = base_buf.GetChapterResponse{
	Status: base_buf.RequestStatus_REQUEST_STATUS_LOADING,
	Data:   nil,
}

var err_chapter base_buf.GetChapterResponse = base_buf.GetChapterResponse{
	Status: base_buf.RequestStatus_REQUEST_STATUS_ERROR,
	Data:   nil,
}

const CHAPTER_KEY string = "CHAPTER"

var CHAPTER_IDEMPOTENCY_KEY string = fmt.Sprintf("%v/IDEMPOTENCY_KEY", CHAPTER_KEY)

type ChapterCache struct {
	Status base_buf.RequestStatus
	Data   *base_buf.Chapter
}

func (handler ChapterHandler) GetChapter(r_ctx context.Context, request *base_buf.GetChapterRequest, stream *connect.ServerStream[base_buf.GetChapterResponse]) error {
	service := usecase.GetMangaChapterService{}
	ctx := *handler.HandlerContext

	redis_ch, ok := ctx.Value(constants.REDIS_URL).(*redis.Client)
	if !ok {
		return errors.New("Invalid redis client")
	}

	err := stream.Send(&initial_chapter)
	if err != nil {
		log.Printf("[CHAPTER] - Unable to send initial status, error %+v\n", err)
		return err
	}

	idempotency_key := fmt.Sprintf("%v/%v/%v/%v", CHAPTER_IDEMPOTENCY_KEY, request.Slug, request.Chapter, request.Id)
	idempotency, err := redis_ch.Get(r_ctx, idempotency_key).Result()
	if err == nil {
		// TODO: Abstract into a separate function
		cache := ChapterCache{}

		cache_err := json.Unmarshal([]byte(idempotency), &cache)
		if cache_err == nil {
			err := stream.Send(&base_buf.GetChapterResponse{
				Status: cache.Status,
				Data:   cache.Data,
			})

			if err != nil {
				log.Printf("[CHAPTER] - Unable to send idempotency cache, error %+v\n", err)
				return err
			}

			return nil
		}

		log.Printf("[CHAPTER] - Unable to parse manga cache, error %+v\n", err)
	}

	amq_connection, ok := ctx.Value(broker.BROKER_CONNECTION).(*amqp091.Connection)
	if !ok {
		log.Printf("[CHAPTER] - AMQ Connection not found for %v\n", events.SCRAPE_CHAPTER_EVENT)
		return errors.New("Invalid amq connection")
	}

	channel, err := amq_connection.Channel()
	if err != nil {
		log.Printf("[CHAPTER] - AMQ Channel not found for %v\n", events.SCRAPE_CHAPTER_EVENT)
		return errors.Join(errors.New("Invalid queue publisher"), err)
	}

	queue_publisher_raw, ok := ctx.Value(events.SCRAPE_CHAPTER_EVENT).(*events.Publisher)
	if !ok {
		log.Printf("[CHAPTER] - Queue (%v) wasn't found\n", events.SCRAPE_CHAPTER_EVENT)
		return errors.New("Invalid queue publisher")
	}

	queue_publisher := *queue_publisher_raw
	service_request := usecase.GetMangaChapterRequest{
		Id:      request.Id,
		Slug:    request.Slug,
		Chapter: request.Chapter,

		Channel:   channel,
		Publisher: queue_publisher,
	}

	service_response, err := service.Exec(service_request)
	if err != nil {
		stream_err := stream.Send(&err_chapter)
		if stream_err != nil {
			log.Printf("[CHAPTER] - Unable to send service error status, error %+v\n", err)
		}

		return err
	}

	chapter := service_response.Chapter.ToProtobuf()
	response := base_buf.GetChapterResponse{
		Data:   &chapter,
		Status: base_buf.RequestStatus_REQUEST_STATUS_COMPLETED,
	}

	cache_idempotency, err := json.Marshal(ChapterCache{
		Status: response.Status,
		Data:   response.Data,
	})
	if err == nil {
		_, err = redis_ch.Set(r_ctx, idempotency_key, cache_idempotency, time.Duration(channels.QUEUE_TTL)*time.Second).Result()

		if err != nil {
			log.Printf("[CHAPTER] - Unable to write idempotency cache, error: %v", err)
		}
	}

	err = stream.Send(&response)
	if err != nil {
		log.Printf("[CHAPTER] - Unable to send service response, error %+v\n", err)
		return err
	}

	return nil
}

func SetupChapterRoute(mux *http.ServeMux, ctx *context.Context) {
	service := ChapterHandler{
		HandlerContext: ctx,
	}

	path, handler := buf_handler.NewChapterServiceHandler(service)
	mux.Handle(path, utils.Middlewares(handler, ctx))
}
