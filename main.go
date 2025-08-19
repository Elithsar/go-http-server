package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(res, req)
	})
}

func (cfg *apiConfig) serverHitsHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(fmt.Sprintf(`
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) serverHitsResetHandler(res http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
	res.Write([]byte("Hits reset."))
}

func readinessHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(200)
	res.Write([]byte("OK"))
}

func main() {
	serveMux := http.NewServeMux()
	cfg := apiConfig{}
	server := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	serveMux.Handle("/app/", cfg.middlewareMetricsInc(handler))
	serveMux.HandleFunc("GET /api/healthz", readinessHandler)
	serveMux.HandleFunc("GET /admin/metrics", cfg.serverHitsHandler)
	serveMux.HandleFunc("POST /admin/reset", cfg.serverHitsResetHandler)

	err := http.ListenAndServe(server.Addr, server.Handler)

	if err != nil {
		fmt.Errorf("failed to init server: %w", err)
	}
}
