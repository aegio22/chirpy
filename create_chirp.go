package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/aegio22/chirpy/internal/auth"
	"github.com/aegio22/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	Id           uuid.UUID     `json:"id"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	Body         string        `json:"body"`
	UserId       uuid.NullUUID `json:"user_id"`
	responseBody struct {
		CleanedBody string `json:"cleaned_body"`
	}
}

func (cfg *apiConfig) handlerCreateChirp(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	chirp := Chirp{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&chirp)
	if err != nil {
		log.Printf("error decoding body into an email: %v", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	authHeader := req.Header

	bearerToken, err := auth.GetBearerToken(authHeader)
	if err != nil {
		log.Printf("error getting bearer token: %v", err)
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	userId, err := auth.ValidateJWT(bearerToken, cfg.jwtSecret)
	if err != nil {
		log.Printf("error validating jwt: %v", err)
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	dbUserID := uuid.NullUUID{
		UUID:  userId,
		Valid: true,
	}

	//start validation logic
	if len(chirp.Body) > 140 {
		lengthErr, _ := json.Marshal(ValidationErr{
			Error: "Chirp is too long",
		})
		rw.WriteHeader(400)
		rw.Write(lengthErr)

		return
	}
	//un-comment this if you want simply to return validity value, as well as switching the rw.Write statement below.
	//validityResponse, err := json.Marshal(ChirpValidity{Valid: true})
	chirp.replaceProfanity()
	chirp.Body = chirp.responseBody.CleanedBody

	queriedChirp, err := cfg.db.CreateChirp(ctx, database.CreateChirpParams{
		Body:   chirp.Body,
		UserID: dbUserID,
	})
	chirp.Id = queriedChirp.ID
	chirp.CreatedAt = queriedChirp.CreatedAt
	chirp.UpdatedAt = queriedChirp.UpdatedAt
	chirp.UserId = dbUserID

	if err != nil {
		log.Fatalf("error adding chirp to database: %v", err)
		return
	}
	chirpJson, err := json.Marshal(chirp)
	if err != nil {
		valErr, _ := json.Marshal(ValidationErr{
			Error: "Error marshalling validated chirp",
		})

		rw.WriteHeader(400)
		rw.Write(valErr)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(201)
	rw.Write(chirpJson)

}
