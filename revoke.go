package main

import (
	"log"
	"net/http"

	"github.com/aegio22/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("error getting token from request: %v", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	err = cfg.db.RevokeToken(ctx, bearerToken)
	if err != nil {
		log.Printf("error revoking token: %v", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusNoContent)
}
