package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/matthieukhl/chirpy/internal/auth"
)

func (a *ApiConfig) PostLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	params := new(parameters)
	err := json.NewDecoder(r.Body).Decode(params)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusBadRequest)
		return
	}

	if params.Email == "" || params.Password == "" {
		log.Println("email and password cannot be empty")
		respondWithError(w, http.StatusBadRequest)
		return
	}

	if params.ExpiresInSeconds == 0 || params.ExpiresInSeconds > 3600 {
		params.ExpiresInSeconds = 3600
	}

	user, err := a.Queries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Incorrect email or password"})
			return
		}
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Incorrect email or password"})
		return
	}

	isPasswordMatch, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Incorrect email or password"})
		return
	}

	if !isPasswordMatch {
		log.Println("Passwords don't match")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Incorrect email or password"})
		return
	}

	expiresIn := time.Duration(params.ExpiresInSeconds) * time.Second

	token, err := auth.MakeJWT(user.ID, a.JWTSecret, expiresIn)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"id":         user.ID.String(),
		"email":      user.Email,
		"created_at": user.CreatedAt.String(),
		"updated_at": user.UpdatedAt.String(),
		"token":      token,
	})

}
