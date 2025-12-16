package main

import (
	"log"
	"net/http"
)

func main() {
	multiPlexer := http.NewServeMux()
	const rootDir = http.Dir(".")
	multiPlexer.Handle("/app/", http.StripPrefix("/app/", http.FileServer(rootDir)))
	multiPlexer.HandleFunc("/healthz", handlerReadiness)
	server := http.Server{
		Addr:    ":8080",
		Handler: multiPlexer,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("error starting server: %v", err)
	}

}

func handlerReadiness(rw http.ResponseWriter, req *http.Request) {

	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("OK"))

}
