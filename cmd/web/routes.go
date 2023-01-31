package main

import (
	"github.com/bmizerany/pat"
	"net/http"
	"websockets/internal/handlers"
)

// routes defines the application routes
func routes() http.Handler {
	mux := pat.New()

	mux.Get("/", http.HandlerFunc(handlers.Home))

	return mux
}
