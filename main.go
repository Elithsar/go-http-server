package main

import(
    "fmt"
    "net/http"
    "sync/atomic"
)

type apiConfig struct {
    fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
    cfg.fileserverHits.Store(cfg.fileserverHits.Add(1))
    return next
}

func readinessHandler(res http.ResponseWriter, req *http.Request) {
    res.Header().Set("Content-Type", "text/plain; charset=utf-8")
    res.WriteHeader(200)
    res.Write([]byte("OK"))
}

func (cfg *apiConfig) serverHitsHandler(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte(fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) serverHitsResetHandler(res http.ResponseWriter, req *http.Request) {
    cfg.fileserverHits.Store(0)
    res.Write([]byte(fmt.Sprintf("Hits reset."))) 
}

func main(){
    serveMux := http.NewServeMux()
    cfg := apiConfig{}
    server := http.Server{
	Addr: ":8080",
	Handler: serveMux,
    }
    handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
    serveMux.Handle("/app/", cfg.middlewareMetricsInc(handler))
    serveMux.HandleFunc("/healthz", readinessHandler)
    serveMux.HandleFunc("/metrics", cfg.serverHitsHandler)
    serveMux.HandleFunc("/reset", cfg.serverHitsResetHandler)

    err := http.ListenAndServe(server.Addr, server.Handler)

    if err != nil {                                   fmt.Errorf("failed to init server: %w", err)
    }
}
