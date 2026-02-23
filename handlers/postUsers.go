package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/matthieukhl/chirpy/internal/auth"
	"github.com/matthieukhl/chirpy/internal/database"
)

func (a *ApiConfig) PostUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	params := new(parameters)

	err := json.NewDecoder(r.Body).Decode(params)
	if err != nil {
		a.Logger.Error(err.Error(), "endpoint", "POST /api/users")
		respondWithError(w, http.StatusBadRequest)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)

	args := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	}

	user, err := a.Queries.CreateUser(r.Context(), args)
	if err != nil {
		a.Logger.Error(err.Error(), "endpoint", "POST /api/users")
		respondWithError(w, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
