package grpc

import (
	"context"
	"errors"
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
	l := utils.GetLogger()
	ctx := *service.ServiceContext

	queue_publisher_raw, ok := ctx.Value(events.SEARCH_MANGA_EVENT).(*events.Publisher)
	if !ok {
		l.Errorf("[PUBLISHER] - Queue (%v) wasn't found", events.SEARCH_MANGA_EVENT)
		return errors.New("Invalid queue publisher")
	}

	queue_publisher := *queue_publisher_raw

	response := base_buf.SearchResponse{
		Status: base_buf.RequestStatus_REQUEST_STATUS_LOADING,
		Data:   nil,
	}

	err := stream.Send(&response)
	if err != nil {
		l.Errorf("[PUBLISHER] - Unable to send initial status, error %+v", err)
		return err
	}

	message := events.SearchMangaMessage{
		Query: request.Slug,
		BaseEvent: events.BaseEvent{
			Metadata: events.Metadata{
				Id: request.Id,
			},
		},
	}

	// Send async message and wait for REDIS update.
	queue_publisher(message)

	// Retrieve REDIS update every 5-10s
	// TODO: implement REDIS conn

	response = base_buf.SearchResponse{
		Status: base_buf.RequestStatus_REQUEST_STATUS_COMPLETED,
		Data:   nil,
	}

	err = stream.Send(&response)
	if err != nil {
		l.Errorf("[PUBLISHER] - Unable to send status update, error %+v", err)
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
