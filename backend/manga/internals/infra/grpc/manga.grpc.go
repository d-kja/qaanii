package grpc

import (
	"context"
	"log"
	"net/http"
	base_buf "qaanii/mangabuf/gen/manga/v1"
	buf_handler "qaanii/mangabuf/gen/manga/v1/mangav1connect"

	"connectrpc.com/connect"
)

type MangaService struct {
	ServiceContext *context.Context
	buf_handler.MangaServiceHandler
}

func (service MangaService) GetManga(_ context.Context, request *base_buf.GetMangaRequest, stream *connect.ServerStream[base_buf.GetMangaResponse]) error {
	response := base_buf.GetMangaResponse{
		Status: base_buf.RequestStatus_REQUEST_STATUS_LOADING,
		Data:   nil,
	}

	err := stream.Send(&response)
	if err != nil {
		return err
	}

	log.Println("Hey!!!")

	stream.Send(nil)
	return nil
}

func SetupMangaRoute(mux *http.ServeMux, ctx *context.Context) {
	service := MangaService{
		ServiceContext: ctx,
	}

	_, handler := buf_handler.NewMangaServiceHandler(service)
	mux.Handle("/manga", handler)
}
