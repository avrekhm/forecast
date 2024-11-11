package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"weather/internal/nws"

	"github.com/gorilla/mux"
)

func NWSForecast(ctx context.Context, client nws.NWSer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		lat, laterr := strconv.ParseFloat(vars["lat"], 64)
		long, longerr := strconv.ParseFloat(vars["long"], 64)

		// this combines error checking for both parsefloats
		if laterr != nil || longerr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		resp, err := nws.GetForecast(ctx, client, lat, long)
		if err != nil {
			log.Printf("error %v", err)

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		payload, err := json.Marshal(resp)
		if err != nil {
			log.Printf("error %v", err)

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(payload)
	}
}
