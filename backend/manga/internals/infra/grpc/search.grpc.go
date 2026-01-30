package grpc

import (
	"context"
	"net/http"
	base_buf "qaanii/mangabuf/gen/manga/v1"
	buf_handler "qaanii/mangabuf/gen/manga/v1/mangav1connect"

	"connectrpc.com/connect"
)

type SearchService struct {
	ServiceContext *context.Context
	buf_handler.SearchServiceHandler
}

func (service SearchService) Search(_ context.Context, request *base_buf.SearchRequest, stream *connect.ServerStream[base_buf.SearchResponse]) error {
	response := base_buf.SearchResponse{
		Status: base_buf.RequestStatus_REQUEST_STATUS_LOADING,
		Data:   nil,
	}

	err := stream.Send(&response)
	if err != nil {
		return err
	}

	return nil
}

func SetupSearchRoute(mux *http.ServeMux, ctx *context.Context) {
	service := SearchService{
		ServiceContext: ctx,
	}

	path, handler := buf_handler.NewSearchServiceHandler(service)
	mux.Handle(path, handler)
}
