package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (a *ApiConfig) GetChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		a.Logger.Error(err.Error(), "endpoint", "GET /api/chirps/{id}")
		respondWithError(w, http.StatusBadRequest)
		return
	}

	chirp, err := a.Queries.GetChirpByID(r.Context(), chirpUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			a.Logger.Warn(err.Error(), "endpoint", "GET /api/chirps/{id}")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		a.Logger.Error(err.Error(), "endpoint", "GET /api/chirps/{id}")
		respondWithError(w, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chirp)
}
