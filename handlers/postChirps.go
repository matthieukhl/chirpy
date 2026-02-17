package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/matthieukhl/chirpy/internal/database"
)

func (a *ApiConfig) PostChirp(w http.ResponseWriter, r *http.Request) {
	type payload struct {
		Body   string `json:"body"`
		UserID string `json:"user_id"`
	}

	requestBody := new(payload)
	err := json.NewDecoder(r.Body).Decode(requestBody)
	if err != nil {
		log.Printf("failed to decode payload: %v", err)
		respondWithError(w, http.StatusBadRequest)
		return
	}

	// Check if user_id is a uuid
	userUUID, err := uuid.Parse(requestBody.UserID)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusBadRequest)
		return
	}

	// Check chirp length
	if len(requestBody.Body) > 140 {
		log.Println(err)
		respondWithError(w, http.StatusBadRequest)
		return
	}

	// Sanitize chirp
	sliceBody := strings.Split(requestBody.Body, " ")
	for i := range sliceBody {
		if strings.ToLower(sliceBody[i]) == "kerfuffle" || strings.ToLower(sliceBody[i]) == "sharbert" || strings.ToLower(sliceBody[i]) == "fornax" {
			sliceBody[i] = "****"
		}
	}

	cleanedBody := strings.Join(sliceBody, " ")

	args := database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: userUUID,
	}

	chirp, err := a.Queries.CreateChirp(r.Context(), args)
	if err != nil {
		log.Printf("failed to create chirp: %v", err)
		respondWithError(w, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(chirp)

}
