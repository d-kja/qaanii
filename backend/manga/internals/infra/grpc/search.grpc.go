package grpc

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"qaanii/manga/internals/infra/broker"
	"qaanii/manga/internals/utils"
	base_buf "qaanii/mangabuf/gen/manga/v1"
	mangav1 "qaanii/mangabuf/gen/manga/v1"
	buf_handler "qaanii/mangabuf/gen/manga/v1/mangav1connect"
	"qaanii/shared/broker/channels"
	"qaanii/shared/broker/events"
	"time"

	"connectrpc.com/connect"
	"github.com/rabbitmq/amqp091-go"
)

type SearchService struct {
	ServiceContext *context.Context
	buf_handler.SearchServiceHandler
}

var err_response base_buf.SearchResponse = base_buf.SearchResponse{
	Status: base_buf.RequestStatus_REQUEST_STATUS_ERROR,
	Data:   nil,
}

func (service SearchService) Search(_ context.Context, request *base_buf.SearchRequest, stream *connect.ServerStream[base_buf.SearchResponse]) error {
	ctx := *service.ServiceContext

	channel, ok := ctx.Value(broker.BROKER_CHANNEL).(*amqp091.Channel)
	if !ok {
		log.Printf("[PUBLISHER] - AMQ Channel not found for %v\n", events.SEARCH_MANGA_EVENT)
		return errors.New("Invalid queue publisher")
	}

	queue_publisher_raw, ok := ctx.Value(events.SEARCH_MANGA_EVENT).(*events.Publisher)
	if !ok {
		log.Printf("[PUBLISHER] - Queue (%v) wasn't found\n", events.SEARCH_MANGA_EVENT)
		return errors.New("Invalid queue publisher")
	}

	queue_publisher := *queue_publisher_raw

	response := base_buf.SearchResponse{
		Status: base_buf.RequestStatus_REQUEST_STATUS_LOADING,
		Data:   nil,
	}

	err := stream.Send(&response)
	if err != nil {
		log.Printf("[PUBLISHER] - Unable to send initial status, error %+v\n", err)
		return err
	}

	// TODO: Search with REDIS for idempotency key/id

	// If not found, send message to queue for async work
	reply_queue, err := channels.CreateReplyQueue(channel)
	if err != nil {
		log.Printf("[PUBLISHER] - Unable to create reply queue, error %+v\n", err)
		return err
	}

	message := events.SearchMangaMessage{
		Query: request.Slug,
		BaseEvent: events.BaseEvent{
			Metadata: events.Metadata{
				Id:    request.Id,
				Reply: reply_queue.Name,
			},
		},
	}

	queue_publisher(message)

	reply_message, err := channel.Consume(reply_queue.Name, "", true, true, false, false, nil)
	if err != nil {
		log.Printf("[PUBLISHER] - Unable to consume reply queue messages, error %+v\n", err)
		return err
	}

	select {

	case reply := <-reply_message:
		{
			message := events.SearchedMangaMessage{}
			err := json.Unmarshal(reply.Body, &message)
			if err != nil {
				log.Printf("[PUBLISHER] - Unable to parse reply body, error %+v\n", err)

				err = stream.Send(&err_response)
				if err != nil {
					log.Printf("[PUBLISHER] - Unable to send status update, error %+v\n", err)
					return err
				}

				return err
			}

			mangas := []*mangav1.Manga{}
			for _, manga := range message.Data {
				buf_manga := manga.ToProtobuf()
				mangas = append(mangas, &buf_manga)
			}

			response = base_buf.SearchResponse{
				Status: base_buf.RequestStatus_REQUEST_STATUS_COMPLETED,
				Data:   mangas,
			}

			err = stream.Send(&response)
			if err != nil {
				log.Printf("[PUBLISHER] - Unable to send status update, error %+v\n", err)
				return err
			}

			break
		}

	case <-time.After(60 * time.Second):
		{
			err = stream.Send(&err_response)
			if err != nil {
				log.Printf("[PUBLISHER] - Unable to send status update, error %+v\n", err)
				return err
			}

			break
		}
	}

	return nil
}

func SetupSearchRoute(mux *http.ServeMux, ctx *context.Context) {
	service := SearchService{
		ServiceContext: ctx,
	}

	path, handler := buf_handler.NewSearchServiceHandler(service)
	mux.Handle(path, utils.Middlewares(handler, ctx))
}
