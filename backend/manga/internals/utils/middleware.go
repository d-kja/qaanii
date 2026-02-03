package utils

import (
	"context"
	"log"
	"net/http"
	"time"
)

func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		startTime := time.Now()

		next.ServeHTTP(writer, req)

		elapsedTime := time.Since(startTime)
		log.Printf("[%s] [%s] [%s]\n", req.Method, req.URL.Path, elapsedTime)
	})
}

func Middlewares(next http.Handler, ctx *context.Context) http.Handler {
	l := GetLogger()

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		url := request.URL
		method := request.Method

		l.Infof("%v request to %v", method, url)

		next.ServeHTTP(writer, request)
	})
}
