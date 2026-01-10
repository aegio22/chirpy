package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aegio22/chirpy/internal/auth"
	"github.com/google/uuid"
)

type WebHookRequest struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) handlerUpgradeUserToChirpyRed(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	webhookRequest := WebHookRequest{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&webhookRequest)
	if err != nil {
		log.Printf("error decoding request body: %v", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if webhookRequest.Event != "user.upgraded" {
		rw.WriteHeader(http.StatusNoContent)
		return
	}

	userId, err := uuid.Parse(webhookRequest.Data.UserID)
	if err != nil {
		log.Printf("error getting user id from request: %v", err)
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		log.Printf("error getting api key: %v", err)
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	if apiKey != cfg.polkaKey {
		log.Printf("api key does not match credentials")
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = cfg.db.UpgradeUserToChirpyRed(ctx, userId)
	if err != nil {
		log.Printf("error updating membership info: %v", err)
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusNoContent)

}
