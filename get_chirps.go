package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"

	"github.com/aegio22/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	// Check if filtering by author_id
	authorId := req.URL.Query().Get("author_id")
	sortParam := req.URL.Query().Get("sort")

	if authorId != "" {
		authUuid, err := uuid.Parse(authorId)
		if err != nil {
			log.Printf("error parsing author id: %v", err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		chirpsQuery, err := cfg.db.GetChirpByAuthorID(ctx, uuid.NullUUID{
			UUID:  authUuid,
			Valid: true,
		})
		if err != nil {
			log.Printf("error getting chirps by author: %v", err)
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
		rw.WriteHeader(http.StatusOK)
		rw.Write(responseBody)
		return
	}

	// No filter - get all chirps
	chirpsQuery, err := cfg.db.GetChirps(ctx)
	if err != nil {
		log.Printf("error getting chirps: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	if sortParam != "" {
		switch sortParam {
		case "asc":
			sort.Slice(chirpsQuery, func(i, j int) bool {
				return i < j
			})
		case "desc":
			sort.Slice(chirpsQuery, func(i, j int) bool {
				return i > j
			})
		default:
			break
		}

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
	rw.WriteHeader(http.StatusOK)
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
