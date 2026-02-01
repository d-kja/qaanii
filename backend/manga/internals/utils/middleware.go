package utils

import (
	"log"
	"net/http"
	"time"
)

type Middleware func(next http.Handler) http.Handler

func Chain(middleware ...Middleware) http.Handler {
	var handler http.Handler
	for i := range middleware {
		handler = middleware[len(middleware)-1-i](handler)
	}

	return handler
}

func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		startTime := time.Now()

		next.ServeHTTP(writer, req)

		elapsedTime := time.Since(startTime)
		log.Printf("[%s] [%s] [%s]\n", req.Method, req.URL.Path, elapsedTime)
	})
}
