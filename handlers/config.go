package handlers

import (
	"database/sql"
	"log"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	"github.com/matthieukhl/chirpy/internal/database"
)

type ApiConfig struct {
	FileServerHits atomic.Int32
	Queries        *database.Queries
}

func NewDB() *database.Queries {
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
