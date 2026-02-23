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
	mux.HandleFunc("POST /api/chirps", apiCfg.PostChirp)

	// Create a new user
	mux.HandleFunc("POST /api/users", apiCfg.PostUsers)

	// GET all chirps
	mux.HandleFunc("GET /api/chirps", apiCfg.GetChirps)

	// GET chirp by ID
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.GetChirpByID)

	// POST login
	mux.HandleFunc("POST /api/login", apiCfg.PostLogin)

	mux.HandleFunc("POST /api/revoke", apiCfg.PostRevoke)

	mux.HandleFunc("POST /api/refresh", apiCfg.PostRefresh)

	// PUT /api/users
	mux.HandleFunc("PUT /api/users", apiCfg.PutUsers)

	// DELETE /api/chirps
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.DeleteChirps)

	// POST /api/polka/webhooks
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.PostPolkaWebhooks)
	// Serve HTTP
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
}
