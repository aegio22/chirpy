package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type UserEmail struct {
	Email string `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(rw http.ResponseWriter, req *http.Request) {
	//create context.Context for SQLC queries. Create user email for JSON decoding
	ctx := req.Context()
	email := UserEmail{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&email)
	if err != nil {
		log.Printf("error decoding body into an email: %v", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	queriedUser, err := cfg.db.CreateUser(ctx, email.Email)
	if err != nil {
		log.Printf("error creating user: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	user := User{
		ID:        queriedUser.ID,
		CreatedAt: queriedUser.CreatedAt,
		UpdatedAt: queriedUser.UpdatedAt,
		Email:     queriedUser.Email,
	}

	responseBody, err := json.Marshal(user)
	if err != nil {
		log.Printf("error marshaling response body: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	rw.Write(responseBody)

}
