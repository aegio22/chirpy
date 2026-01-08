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
	
	ctx := req.Context()
	
	// Delete refresh tokens first to avoid foreign key constraint violation
	err := cfg.db.DeleteAllRefreshTokens(ctx)
	if err != nil {
		log.Printf("error deleting refresh tokens: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	// Now delete users
	err = cfg.db.DeleteUsers(ctx)
	if err != nil {
		log.Printf("error deleting users: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	cfg.fileserverHits.Store(0)
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("OK"))
}
