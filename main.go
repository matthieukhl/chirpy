package main

import (
	"fmt"
	"log"
	"net/http"
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

	// Serve HTTP
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
}
