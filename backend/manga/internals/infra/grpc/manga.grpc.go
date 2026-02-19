package grpc

import (
	"context"
	"errors"
	"log"
	"net/http"
	"qaanii/manga/internals/constants"
	usecase "qaanii/manga/internals/domain/mangas/use_case"
	"qaanii/manga/internals/utils"
	base_buf "qaanii/mangabuf/gen/manga/v1"
	buf_handler "qaanii/mangabuf/gen/manga/v1/mangav1connect"

	"connectrpc.com/connect"
	"github.com/redis/go-redis/v9"
)

var initial_status base_buf.GetMangaResponse = base_buf.GetMangaResponse{
	Status: base_buf.RequestStatus_REQUEST_STATUS_LOADING,
	Data:   nil,
}

var err_response base_buf.GetMangaResponse = base_buf.GetMangaResponse{
	Status: base_buf.RequestStatus_REQUEST_STATUS_ERROR,
	Data:   nil,
}

type MangaService struct {
	HandlerContext *context.Context
	buf_handler.MangaServiceHandler
}

func (handler MangaService) GetManga(r_ctx context.Context, request *base_buf.GetMangaRequest, stream *connect.ServerStream[base_buf.GetMangaResponse]) error {
	service := usecase.GetMangaBySlugService{}
	ctx := *handler.HandlerContext

	redis_ch, ok := ctx.Value(constants.REDIS_URL).(*redis.Client)
	if !ok {
		return errors.New("Invalid redis client")
	}

	err := stream.Send(&initial_status)
	if err != nil {
		log.Printf("[MANGA] - Unable to send initial status, error %+v\n", err)
		return err
	}

	log.Println("Hey!!!")

	stream.Send(nil)
	return nil
}

func SetupMangaRoute(mux *http.ServeMux, ctx *context.Context) {
	service := MangaService{
		HandlerContext: ctx,
	}

	path, handler := buf_handler.NewMangaServiceHandler(service)
	mux.Handle(path, utils.Middlewares(handler, ctx))
}
