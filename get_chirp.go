package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)



func (cfg *apiConfig) handlerGetChirpByID(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	id := req.PathValue("chirpID")


	chirpID, err := uuid.Parse(id)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	queryChirp, err := cfg.db.GetChirpByID(ctx, chirpID)
	if err != nil {
		log.Printf("error getting chirp by ID from db: %v", err)
		rw.WriteHeader(404)
		return
	}
	chirp := helperSQLChirpToChirp(queryChirp)
	responseBody, err := json.Marshal(chirp)
	if err != nil {
		log.Printf("error marshaling chirp: %v", err)
		rw.WriteHeader(404)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(200)
	rw.Write(responseBody)

}
