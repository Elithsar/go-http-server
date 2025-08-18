package main

import(
    "fmt"
    "net/http"
)

func main(){
    serveMux := http.NewServeMux()
    server := http.Server{
	Addr: ":8080",
	Handler: serveMux,
    }
    serveMux.Handle("/", http.FileServer(http.Dir(".")))

    err := http.ListenAndServe(server.Addr, server.Handler)

    if err != nil {                                   fmt.Errorf("failed to init server: %w", err)
    }
}
