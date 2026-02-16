package handlers

import "net/http"

func (a *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.FileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
