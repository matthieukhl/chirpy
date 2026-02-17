package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func (a *ApiConfig) GetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := a.Queries.GetAllChirps(r.Context())
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(chirps); err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError)
		return
	}
}
