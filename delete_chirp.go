package main

import (
	"log"
	"net/http"

	"github.com/aegio22/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirpByID(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	id := req.PathValue("chirpID")
	chirpID, err := uuid.Parse(id)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	accessToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("error getting access token: %v", err)
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}
	userId, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
	if err != nil {
		log.Printf("error validating token: %v", err)
		rw.WriteHeader(http.StatusForbidden)
		return
	}
	chirp, err := cfg.db.GetChirpByID(ctx, chirpID)
	if err != nil {
		log.Printf("error fetching chirp: %v", err)
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	if userId != chirp.UserID.UUID {
		log.Print("user id does not match the chirps author")
		rw.WriteHeader(http.StatusForbidden)
		return
	}
	err = cfg.db.DeleteChirpByID(ctx, chirp.ID)
	if err != nil {
		log.Printf("error deleting chirp: %v", err)
		rw.WriteHeader(http.StatusBadRequest)

		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(204)

}
