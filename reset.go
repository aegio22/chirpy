package main

import "net/http"

func (cfg *apiConfig) handlerResetHitCount(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("OK"))
	cfg.fileserverHits.Store(0)
}
