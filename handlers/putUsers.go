package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/matthieukhl/chirpy/internal/auth"
	"github.com/matthieukhl/chirpy/internal/database"
)

func (a *ApiConfig) PutUsers(w http.ResponseWriter, r *http.Request) {
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

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Retrieve data from request body
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		a.Logger.Error(err.Error(), "endpoint", r.URL.String())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "bad request"})
		return
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		a.Logger.Error(err.Error(), "endpoint", r.URL.String())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	// Build query arguments
	args := database.UpdateUserInfoParams{
		Email:          req.Email,
		HashedPassword: hashedPassword,
		ID:             userID,
	}

	// Update user information in database
	user, err := a.Queries.UpdateUserInfo(r.Context(), args)
	if err != nil {
		a.Logger.Error(err.Error(), "endpoint", r.URL.String())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	// Return updated user information
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)

}
