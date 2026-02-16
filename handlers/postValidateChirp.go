package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (a *ApiConfig) PostValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest)
		return
	}

	sliceBody := strings.Split(params.Body, " ")
	for i := range sliceBody {
		if strings.ToLower(sliceBody[i]) == "kerfuffle" || strings.ToLower(sliceBody[i]) == "sharbert" || strings.ToLower(sliceBody[i]) == "fornax" {
			sliceBody[i] = "****"
		}
	}

	cleanedBody := strings.Join(sliceBody, " ")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]string{"cleaned_body": cleanedBody})

}
