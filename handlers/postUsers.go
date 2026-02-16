package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func (a *ApiConfig) PostUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	params := new(parameters)

	err := json.NewDecoder(r.Body).Decode(params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest)
		return
	}

	user, err := a.Queries.CreateUser(r.Context(), params.Email)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
