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

type User struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"hashed_password"`
	Token          string    `json:"token"`
	RefreshToken   string    `json:"refresh_token"`
	IsChirpyRed    bool      `json:"is_chirpy_red"`
}

type UserInfo struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds *int   `json:"expires_in_seconds"`
}

func (cfg *apiConfig) handlerCreateUser(rw http.ResponseWriter, req *http.Request) {
	//create context.Context for SQLC queries. Create user email for JSON decoding
	ctx := req.Context()
	userInfo := UserInfo{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&userInfo)
	if err != nil {
		log.Printf("error decoding body into UserInfo: %v", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	hashedPassword, err := auth.HashPassword(userInfo.Password)
	if err != nil {
		log.Print(err)
		return
	}
	queriedUser, err := cfg.db.CreateUser(ctx, database.CreateUserParams{Email: userInfo.Email, HashedPassword: hashedPassword})
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
