package main

import(
    "fmt"
    "net/http"
)

func readinessHandler(res http.ResponseWriter, req *http.Request) {
    res.Header().Set("Content-Type", "text/plain; charset=utf-8")
    res.WriteHeader(200)
    res.Write([]byte("OK"))
}

func main(){
    serveMux := http.NewServeMux()
    server := http.Server{
	Addr: ":8080",
	Handler: serveMux,
    }
    serveMux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))	
    serveMux.HandleFunc("/healthz", readinessHandler)

    err := http.ListenAndServe(server.Addr, server.Handler)

    if err != nil {                                   fmt.Errorf("failed to init server: %w", err)
    }
}
