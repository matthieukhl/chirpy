package handlers

import "sync/atomic"

type ApiConfig struct {
	FileServerHits atomic.Int32
}
