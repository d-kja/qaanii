package grpc

import (
	"context"
	"errors"
	"log"
	"net/http"
	"qaanii/manga/internals/utils"
	base_buf "qaanii/mangabuf/gen/manga/v1"
	buf_handler "qaanii/mangabuf/gen/manga/v1/mangav1connect"
	"qaanii/shared/broker/events"

	"connectrpc.com/connect"
)

type SearchService struct {
	ServiceContext *context.Context
	buf_handler.SearchServiceHandler
}

func (service SearchService) Search(_ context.Context, request *base_buf.SearchRequest, stream *connect.ServerStream[base_buf.SearchResponse]) error {
	ctx := *service.ServiceContext

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
	message := events.SearchMangaMessage{
		Query: request.Slug,
		BaseEvent: events.BaseEvent{
			Metadata: events.Metadata{
				Id: request.Id,
			},
		},
	}

	queue_publisher(message)

	// TODO: Wait for queue to reply and update status

	// Mock example
	response = base_buf.SearchResponse{
		Status: base_buf.RequestStatus_REQUEST_STATUS_COMPLETED,
		Data:   nil,
	}

	err = stream.Send(&response)
	if err != nil {
		log.Printf("[PUBLISHER] - Unable to send status update, error %+v\n", err)
		return err
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
