package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/matthieukhl/chirpy/internal/auth"
)

func (a *ApiConfig) PostRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		if err.Error() == "Missing Authorization header" {
			a.Logger.Error(err.Error(), "endpoint", "POST /api/revoke")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
			return
		}
		a.Logger.Error(err.Error(), "endpoint", "POST /api/revoke")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	_, err = a.Queries.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		if err == sql.ErrNoRows {
			a.Logger.Warn(err.Error(), "endpoint", "POST /api/revoke")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "token not found"})
			return
		}
		a.Logger.Error(err.Error(), "endpoint", "POST /api/revoke")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
