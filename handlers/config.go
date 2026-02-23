package handlers

import (
	"database/sql"
	"log"
	"log/slog"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	"github.com/matthieukhl/chirpy/internal/database"
)

type ApiConfig struct {
	FileServerHits atomic.Int32
	Queries        *database.Queries
	Platform       string
	JWTSecret      string
	PolkaApiKey    string
	Logger         *slog.Logger
}

func NewConfig() ApiConfig {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	return ApiConfig{
		Queries:     newDB(),
		Platform:    os.Getenv("PLATFORM"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		PolkaApiKey: os.Getenv("POLKA_KEY"),
		Logger:      logger,
	}
}

func newDB() *database.Queries {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	dbUrl := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	return database.New(db)
}
