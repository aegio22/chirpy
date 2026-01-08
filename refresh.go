package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aegio22/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("error fetching bearer token: %v", err)
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}
	userFetch, err := cfg.db.GetUserFromRefreshToken(ctx, refreshToken)
	if err != nil {
		log.Printf("error fetching user by refresh token: %v", err)
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}
	type TokenBody struct {
		Token string `json:"token"`
	}
	accessToken, err := auth.MakeJWT(userFetch.ID, cfg.jwtSecret)
	if err != nil {
		log.Printf("error generating jwt access token: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	tb := TokenBody{Token: accessToken}
	tokenBody, err := json.Marshal(tb)
	if err != nil {
		log.Printf("error generating response: %v", err)
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(tokenBody)

}
