package handlers

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int) {
	type resp struct {
		Error string `json:"error"`
	}

	w.Header().Set("Content-Type", "application/json")

	// Case - 400 Bad Request
	if code == http.StatusBadRequest {
		w.WriteHeader(http.StatusBadRequest)
		response := resp{Error: "Chirp is too long"}
		dat, err := json.Marshal(response)
		if err != nil {
			slog.Error(err.Error())
			respondWithError(w, http.StatusInternalServerError)
		}
		w.Write(dat)
		return
	}

	// Case - 500 Internal Server Error
	if code == http.StatusInternalServerError {
		w.WriteHeader(http.StatusInternalServerError)
		response := resp{Error: "Something went wrong"}
		dat, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error marshalling JSON: %v\n", err)
		}
		w.Write(dat)
		return
	}

}
