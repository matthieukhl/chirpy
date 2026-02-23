package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/matthieukhl/chirpy/internal/auth"
	"github.com/matthieukhl/chirpy/internal/database"
)

func (a *ApiConfig) PostChirp(w http.ResponseWriter, r *http.Request) {
	type payload struct {
		Body string `json:"body"`
	}

	token, err := auth.GetBearerToken(r.Header)
	log.Println(token)
	if err != nil {
		a.Logger.Error(err.Error(), "endpoint", "POST /api/chirps")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
		return
	}

	userID, err := auth.ValidateJWT(token, a.JWTSecret)
	if err != nil {
		a.Logger.Error(err.Error(), "endpoint", "POST /api/chirps")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
		return
	}

	requestBody := new(payload)
	err = json.NewDecoder(r.Body).Decode(requestBody)
	if err != nil {
		a.Logger.Error(err.Error(), "endpoint", "POST /api/chirps")
		respondWithError(w, http.StatusBadRequest)
		return
	}

	// Check chirp length
	if len(requestBody.Body) > 140 {
		a.Logger.Error("invalid chirp length")
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
		UserID: userID,
	}

	chirp, err := a.Queries.CreateChirp(r.Context(), args)
	if err != nil {
		a.Logger.Error(err.Error(), "endpoint", "POST /api/chirps")
		respondWithError(w, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(chirp)

}
