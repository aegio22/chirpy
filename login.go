package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/aegio22/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(rw http.ResponseWriter, req *http.Request) {
	userInfo := UserInfo{}
	ctx := req.Context()
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&userInfo)
	if err != nil {
		log.Printf("error decoding body into UserInfo: %v", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Set default expiry if not provided
	if userInfo.ExpiresInSeconds == nil {
		defaultExpiry := 3600
		userInfo.ExpiresInSeconds = &defaultExpiry
	}
	// Cap maximum expiry at 3600 seconds
	if *userInfo.ExpiresInSeconds > 3600 {
		maxExpiry := 3600
		userInfo.ExpiresInSeconds = &maxExpiry
	}

	// Get user from database
	loginQuery, err := cfg.db.GetUserByEmail(ctx, userInfo.Email)
	if err != nil {
		log.Print("Incorrect email or password")
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Verify password BEFORE creating JWT
	match, err := auth.CheckPasswordHash(userInfo.Password, loginQuery.HashedPassword)
	if err != nil {
		log.Print("Incorrect email or password")
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}
	if !match {
		log.Print("Incorrect email or password")
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Create JWT token only after password is verified
	jwtToken, err := auth.MakeJWT(loginQuery.ID, cfg.jwtSecret, time.Duration(*userInfo.ExpiresInSeconds)*time.Second)
	if err != nil {
		log.Printf("error creating JWT: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	user := User{
		ID:        loginQuery.ID,
		CreatedAt: loginQuery.CreatedAt,
		UpdatedAt: loginQuery.UpdatedAt,
		Email:     loginQuery.Email,
		Token:     jwtToken,
	}

	responseBody, err := json.Marshal(user)
	if err != nil {
		log.Printf("error marshaling response body: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(responseBody)
}
