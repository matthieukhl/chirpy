package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"sort"

	"github.com/google/uuid"
)

func (a *ApiConfig) GetChirps(w http.ResponseWriter, r *http.Request) {
	// Retrieve optinal author_id parameter
	author_id := r.URL.Query().Get("author_id")
	sortOrder := r.URL.Query().Get("sort")

	if sortOrder != "asc" && sortOrder != "desc" {
		a.Logger.Error("invalid sort parameter", "endpoint", r.URL.String(), "sort_parameter", sortOrder)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "sort must be either 'asc' or 'desc'"})
		return
	}

	if author_id == "" {

		// Check if sortOrder is 'desc'
		if sortOrder == "desc" {
			chirps, err := a.Queries.GetAllChirps(r.Context())

			if err != nil {
				a.Logger.Error(err.Error(), "endpoint", "GET /api/chirps")
				respondWithError(w, http.StatusInternalServerError)
				return
			}

			sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.After(chirps[j].CreatedAt) })
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(chirps); err != nil {
				a.Logger.Error(err.Error(), "endpoint", "GET /api/chirps")
				respondWithError(w, http.StatusInternalServerError)
				return
			}
			return
		}

		// No sortOrder -> by default the query returns with 'ORDER BY created_at ASC'
		chirps, err := a.Queries.GetAllChirps(r.Context())
		if err != nil {
			a.Logger.Error(err.Error(), "endpoint", "GET /api/chirps")
			respondWithError(w, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(chirps); err != nil {
			a.Logger.Error(err.Error(), "endpoint", "GET /api/chirps")
			respondWithError(w, http.StatusInternalServerError)
			return
		}
		return
	}

	authorUUID, err := uuid.Parse(author_id)
	if err != nil {
		a.Logger.Error(err.Error(), "endpoint", r.URL.String())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "bad request"})
		return
	}

	chirps, err := a.Queries.GetChirpsByUserId(r.Context(), authorUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			a.Logger.Error(err.Error(), "endpoint", r.URL.String())
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "user id not found"})
			return
		}
		a.Logger.Error(err.Error(), "endpoint", r.URL.String())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	if sortOrder == "desc" {
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.After(chirps[j].CreatedAt) })
	}

	w.Header().Set("Content-Type", "appllication/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chirps)
}
