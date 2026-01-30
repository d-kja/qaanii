package grpc

import (
	"context"
	"net/http"
	base_buf "qaanii/mangabuf/gen/manga/v1"
	buf_handler "qaanii/mangabuf/gen/manga/v1/mangav1connect"

	"connectrpc.com/connect"
)

type ChapterService struct {
	ServiceContext *context.Context
	buf_handler.ChapterServiceHandler
}

func (service ChapterService) GetChapter(_ context.Context, request *base_buf.GetChapterRequest, stream *connect.ServerStream[base_buf.GetChapterResponse]) error {
	response := base_buf.GetChapterResponse{
		Status: base_buf.RequestStatus_REQUEST_STATUS_LOADING,
		Data:   nil,
	}

	err := stream.Send(&response)
	if err != nil {
		return err
	}

	return nil
}

func SetupChapterRoute(mux *http.ServeMux, ctx *context.Context) {
	service := ChapterService{
		ServiceContext: ctx,
	}

	path, handler := buf_handler.NewChapterServiceHandler(service)
	mux.Handle(path, handler)
}
