package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/matthieukhl/chirpy/internal/auth"
	"github.com/matthieukhl/chirpy/internal/database"
)

func (a *ApiConfig) PostLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	params := new(parameters)
	err := json.NewDecoder(r.Body).Decode(params)
	if err != nil {
		a.Logger.Error(err.Error(), "endpoint", "POST /api/login")
		respondWithError(w, http.StatusBadRequest)
		return
	}

	if params.Email == "" || params.Password == "" {
		a.Logger.Error("email and password cannot be empty")
		respondWithError(w, http.StatusBadRequest)
		return
	}

	user, err := a.Queries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			a.Logger.Warn(err.Error(), "endpoint", "POST /api/login")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Incorrect email or password"})
			return
		}
		a.Logger.Error(err.Error(), "endpoint", "POST /api/login")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Incorrect email or password"})
		return
	}

	isPasswordMatch, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		a.Logger.Error(err.Error(), "endpoint", "POST /api/login")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Incorrect email or password"})
		return
	}

	if !isPasswordMatch {
		a.Logger.Error("wrong password")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Incorrect email or password"})
		return
	}

	expiresIn := time.Duration(1) * time.Hour

	token, err := auth.MakeJWT(user.ID, a.JWTSecret, expiresIn)
	if err != nil {
		a.Logger.Error(err.Error(), "endpoint", "POST /api/login")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	refreshToken := auth.MakeRefreshToken()

	args := database.CreateRefreshTokenParams{
		Token:  refreshToken,
		UserID: user.ID,
	}

	_, err = a.Queries.CreateRefreshToken(r.Context(), args)
	if err != nil {
		a.Logger.Error(err.Error(), "endpoint", "POST /api/login")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"id":            user.ID.String(),
		"email":         user.Email,
		"created_at":    user.CreatedAt.String(),
		"updated_at":    user.UpdatedAt.String(),
		"token":         token,
		"refresh_token": refreshToken,
		"is_chirpy_red": user.IsChirpyRed,
	})

}
