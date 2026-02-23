package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/matthieukhl/chirpy/internal/auth"
)

func (a *ApiConfig) DeleteChirps(w http.ResponseWriter, r *http.Request) {
	// Retrieve JWT
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		a.Logger.Error(err.Error(), "endpoint", r.URL.String())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
		return
	}

	// Validate JWT
	userID, err := auth.ValidateJWT(token, a.JWTSecret)
	if err != nil {
		a.Logger.Error(err.Error(), "endpoint", r.URL.String())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
		return
	}

	chirpID := r.PathValue("chirpID")

	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		a.Logger.Error(err.Error(), "endpoint", r.URL.String())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "bad request"})
		return
	}

	chirp, err := a.Queries.GetChirpByID(r.Context(), chirpUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			a.Logger.Error(err.Error(), "endpoint", r.URL.String())
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "chirp not found"})
			return
		}
		a.Logger.Error(err.Error(), "endpoint", r.URL.String())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	if userID != chirp.UserID {
		a.Logger.Error("user is not allowed to delete chirp", "endpoint", r.URL.String(), "chirp_id", chirpUUID.String(), "user_id", userID.String())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": "forbidden"})
		return
	}

	err = a.Queries.DeleteChirp(r.Context(), chirpUUID)
	if err != nil {
		a.Logger.Error(err.Error(), "endpoint", r.URL.String(), "chirp_id", chirpUUID.String())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
