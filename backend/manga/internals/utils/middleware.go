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

func Recover() {
	if recovered_err := recover(); recovered_err != nil {
		log.Printf("[MIDDLEWARE] - Recovered from error: %v", recovered_err)
	}
}

func Middlewares(next http.Handler, ctx *context.Context) http.Handler {
	return Log(next)
}
