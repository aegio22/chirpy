package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerReset(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if cfg.platform != "dev" {
		rw.WriteHeader(http.StatusForbidden)
		return
	}
	err := cfg.db.DeleteUsers(req.Context())
	if err != nil {
		log.Printf("error deleting users: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	cfg.fileserverHits.Store(0)
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("OK"))
}
