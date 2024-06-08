package main

import (
	"log"
	"net/http"
)

type Middleware func(http.Handler) http.Handler

type RWWrapper struct {
	responseWriter http.ResponseWriter
	statusCode     int
}

func NewRWWrapper(responseWriter http.ResponseWriter) *RWWrapper {
	return &RWWrapper{responseWriter: responseWriter, statusCode: 0}
}

func (rww *RWWrapper) Header() http.Header {
	return rww.responseWriter.Header()
}

func (rww *RWWrapper) Write(b []byte) (int, error) {
	if rww.statusCode == 0 {
		rww.statusCode = http.StatusOK
	}
	return rww.responseWriter.Write(b)
}

func (rww *RWWrapper) WriteHeader(statusCode int) {
	rww.statusCode = statusCode
	rww.responseWriter.WriteHeader(statusCode)
}

func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	for _, m := range middlewares {
		handler = m(handler)
	}
	return handler
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rww := NewRWWrapper(w)
		next.ServeHTTP(rww, r)
		log.Printf("%d %s %s", rww.statusCode, r.Method, r.URL)
	})
}

func ContentHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
