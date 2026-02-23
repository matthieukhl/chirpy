package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

const event = "user.upgraded"

func (a *ApiConfig) PostPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		a.Logger.Error(err.Error(), "endpoint", r.URL.String())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	if req.Event != event {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID, err := uuid.Parse(req.Data.UserID)
	if err != nil {
		a.Logger.Error(err.Error(), "endpoint", r.URL.String(), "user_id", req.Data.UserID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "bad request"})
		return
	}

	_, err = a.Queries.UpgradeChirpyRed(r.Context(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			a.Logger.Error(err.Error(), "endpoint", r.URL.String(), "user_id", userID.String())
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "user id not found"})
			return
		}
		a.Logger.Error(err.Error(), "endpoint", r.URL.String(), "user_id", userID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
