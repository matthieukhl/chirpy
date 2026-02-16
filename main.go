package main

import (
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/matthieukhl/chirpy/handlers"
)

func main() {
	// Init config
	apiCfg := handlers.NewConfig()

	// Register handlers
	mux := http.NewServeMux()

	// File server handler
	mux.Handle("/app/", apiCfg.MiddlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))

	// Health endpoint
	mux.HandleFunc("GET /api/healthz", apiCfg.GetHealthz)

	// Metrics endpoint
	mux.HandleFunc("GET /admin/metrics", apiCfg.GetAdminMetrics)

	// Reset metrics endpoint
	mux.HandleFunc("POST /admin/reset", apiCfg.PostAdminReset)

	// Validate chirp endpoint
	mux.HandleFunc("POST /api/validate_chirp", apiCfg.PostValidateChirp)

	// Create a new user
	mux.HandleFunc("POST /api/users", apiCfg.PostUsers)

	// Serve HTTP
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
}
