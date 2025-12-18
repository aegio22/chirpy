package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	cfg := apiConfig{}
	multiPlexer := http.NewServeMux()
	const rootDir = http.Dir(".")

	multiPlexer.Handle("/app/",
		http.StripPrefix("/app/",
			cfg.middlewareMetricsInc(http.FileServer(rootDir))))

	multiPlexer.HandleFunc("/healthz", handlerReadiness)
	multiPlexer.HandleFunc("/metrics", cfg.handlerCountReqs)
	multiPlexer.HandleFunc("/reset", cfg.handlerResetHitCount)

	server := http.Server{
		Addr:    ":8080",
		Handler: multiPlexer,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("error starting server: %v", err)
	}

}

func (cfg *apiConfig) handlerCountReqs(rw http.ResponseWriter, req *http.Request) {
	hits := cfg.fileserverHits.Load()
	hitsReport := fmt.Sprintf("Hits: %v", hits)

	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(hitsReport))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(rw, req)
	})

}
