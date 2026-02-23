package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/matthieukhl/chirpy/internal/auth"
)

func (a *ApiConfig) PostRefresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		if err.Error() == "Missing Authorization header" {
			a.Logger.Error(err.Error(), "endpoint", "POST /api/refresh")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
			return
		}
		a.Logger.Error(err.Error(), "endpoint", "POST /api/refresh")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	token, err := a.Queries.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		if err == sql.ErrNoRows {
			a.Logger.Error(err.Error(), "endpoint", "POST /api/refresh")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
			return
		}
		a.Logger.Error(err.Error(), "endpoint", "POST /api/refresh")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	if token.ExpiresAt.Compare(time.Now()) == -1 || token.RevokedAt.Valid == true {
		a.Logger.Warn("token is expired")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
		return
	}

	expiresIn := time.Duration(1) * time.Hour

	newToken, err := auth.MakeJWT(token.UserID, a.JWTSecret, expiresIn)
	if err != nil {
		a.Logger.Error(err.Error(), "endpoint", "POST /api/refresh")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": newToken})
}
