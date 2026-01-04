package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aegio22/chirpy/internal/database"
)

func (cfg *apiConfig) handlerGetChirps(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	chirpsQuery, err := cfg.db.GetChirps(ctx)
	if err != nil {
		log.Printf("error getting chirps: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	chirps := make([]Chirp, 0, len(chirpsQuery))
	for _, sqlChirp := range chirpsQuery {
		chirps = append(chirps, helperSQLChirpToChirp(sqlChirp))
	}
	responseBody, err := json.Marshal(chirps)
	if err != nil {
		log.Printf("error marshalling chirps: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(200)
	rw.Write(responseBody)

}

func helperSQLChirpToChirp(sqlChirp database.Chirp) Chirp {
	chirp := Chirp{
		Id:        sqlChirp.ID,
		CreatedAt: sqlChirp.CreatedAt,
		UpdatedAt: sqlChirp.UpdatedAt,
		Body:      sqlChirp.Body,
		UserId:    sqlChirp.UserID,
	}
	return chirp
}
