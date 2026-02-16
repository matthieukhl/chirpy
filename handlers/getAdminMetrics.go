package handlers

import (
	"fmt"
	"net/http"
)

func (a *ApiConfig) GetAdminMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write(fmt.Appendf(nil, `
		<html>
  			<body>
    			<h1>Welcome, Chirpy Admin</h1>
   	 			<p>Chirpy has been visited %d times!</p>
  			</body>
		</html>`, a.FileServerHits.Load()))
}
