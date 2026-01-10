package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aegio22/chirpy/internal/auth"
	"github.com/aegio22/chirpy/internal/database"
)

func (cfg *apiConfig) handlerChangeUserInfo(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	updatedInfo := UserInfo{}
	decoder := json.NewDecoder(req.Body)
	accessToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("error getting access token: %v", err)
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = decoder.Decode(&updatedInfo)
	if err != nil {
		log.Printf("error decoding request body: %v", err)
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	hashedPassword, err := auth.HashPassword(updatedInfo.Password)
	if err != nil {
		log.Printf("error hashing password: %v", err)
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	userId, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
	if err != nil {
		log.Printf("error validating access token: %v", err)
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := cfg.db.UpdateUserInfo(ctx, database.UpdateUserInfoParams{
		Email:          updatedInfo.Email,
		HashedPassword: hashedPassword,
		ID:             userId,
	})
	if err != nil {
		log.Printf("error updating user info: %v", err)
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	//now fetch updated user

	returnedUser := User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}

	responseBody, err := json.Marshal(returnedUser)
	if err != nil {
		log.Printf("error marshaling response body: %v", err)
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(responseBody)

}
