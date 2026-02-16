package handlers

import (
	"log"
	"net/http"
)

func (a *ApiConfig) PostAdminReset(w http.ResponseWriter, r *http.Request) {
	if a.Platform != "dev" {
		respondWithError(w, http.StatusForbidden)
		return
	}

	err := a.Queries.DeleteAllUsers(r.Context())
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	a.FileServerHits.Store(0)
}
