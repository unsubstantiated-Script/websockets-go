package main

import (
	"log"
	"net/http"
	"websockets/internal/handlers"
)

func main() {
	mux := routes()

	log.Println("Starting Channel Listener")

	//Firing up our go routine that get's this jazz going
	go handlers.ListenToWsChannel()

	log.Println("Starting web server on port 8080")

	_ = http.ListenAndServe(":8080", mux)
}
