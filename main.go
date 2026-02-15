package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func (a *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Init config
	apiCfg := apiConfig{}

	// Register handlers
	mux := http.NewServeMux()

	// File server handler
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))

	// Health endpoint
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Metrics endpoint
	mux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write(fmt.Appendf(nil, `
		<html>
  			<body>
    			<h1>Welcome, Chirpy Admin</h1>
   	 			<p>Chirpy has been visited %d times!</p>
  			</body>
		</html>`, apiCfg.fileServerHits.Load()))
	})

	// Reset metrics endpoint
	mux.HandleFunc("POST /admin/reset", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		apiCfg.fileServerHits.Store(0)
	})

	// Validate chirp endpoint
	mux.HandleFunc("POST /api/validate_chirp", func(w http.ResponseWriter, r *http.Request) {
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

	})

	// Serve HTTP
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
}

func respondWithError(w http.ResponseWriter, code int) {
	type resp struct {
		Error string `json:"error"`
	}

	w.Header().Set("Content-Type", "application/json")

	// Case - 400 Bad Request
	if code == http.StatusBadRequest {
		w.WriteHeader(http.StatusBadRequest)
		response := resp{Error: "Chirp is too long"}
		dat, err := json.Marshal(response)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError)
		}
		w.Write(dat)
		return
	}

	// Case - 500 Internal Server Error
	if code == http.StatusInternalServerError {
		w.WriteHeader(http.StatusInternalServerError)
		response := resp{Error: "Something went wrong"}
		dat, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error marshalling JSON: %v\n", err)
		}
		w.Write(dat)
		return
	}

}
