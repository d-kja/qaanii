package grpc

import (
	"context"
	"net/http"
)

type GRPC struct{}

func (GRPC) Setup(service_context *context.Context) (*http.ServeMux, *http.Protocols) {
	mux := http.NewServeMux()

	// Create routes
	Router(mux, service_context)

	protocol := new(http.Protocols)
	protocol.SetHTTP1(true)
	protocol.SetUnencryptedHTTP2(true) // We are only running locally, we don't need TLS

	return mux, protocol
}

func Router(mux *http.ServeMux, ctx *context.Context) {
	SetupSearchRoute(mux, ctx)
	SetupMangaRoute(mux, ctx)
	SetupChapterRoute(mux, ctx)
}
