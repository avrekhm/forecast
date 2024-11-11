package main

import (
	"context"
	"log"
	"net/http"
	"weather/internal/handlers"
	"weather/internal/nws"

	"github.com/gorilla/mux"
)

var serverPort = "8080"

func main() {
	ctx := context.Background()
	nwsClient := nws.NewClient(ctx)

	router := mux.NewRouter()
	router.HandleFunc("/forecast/{lat},{long}", handlers.NWSForecast(ctx, nwsClient))

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:" + serverPort,
	}

	log.Printf("listening at %s", srv.Addr)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
