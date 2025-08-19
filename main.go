package main

import (
	"encoding/json"
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

func chirpIn(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Body string `json:"body"`
	}
	type response struct {
		Valid bool   `json:"valid,omitempty"`
		Error string `json:"error,omitempty"`
	}
	w.Header().Set("Content-Type", "application/json")
	res := response{}

	decoder := json.NewDecoder(r.Body)
	var message request
	err := decoder.Decode(&message)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res.Error = "Something went wrong"
		out, _ := json.Marshal(res)
		w.Write(out)
		return
	}

	MAX_SIZE := 140
	if len(message.Body) > MAX_SIZE {
		w.WriteHeader(http.StatusBadRequest)
		res.Error = "Chirp is too long"
		out, _ := json.Marshal(res)
		w.Write(out)
		return
	}

	w.WriteHeader(http.StatusOK)

	res.Valid = true
	jsonResp, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong"))
		return
	}
	w.Write(jsonResp)
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
	serveMux.HandleFunc("POST /api/validate_chirp", chirpIn)

	serveMux.HandleFunc("GET /admin/metrics", cfg.serverHitsHandler)
	serveMux.HandleFunc("POST /admin/reset", cfg.serverHitsResetHandler)

	err := http.ListenAndServe(server.Addr, server.Handler)

	if err != nil {
		fmt.Errorf("failed to init server: %w", err)
	}
}
