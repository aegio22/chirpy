package main

import (
	"strings"
)

type ValidationErr struct {
	Error string `json:"error"`
}

type ChirpValidity struct {
	Valid bool `json:"valid"`
}

func (chirp *Chirp) replaceProfanity() {
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}

	chirpWords := strings.Split(chirp.Body, " ")
	for i, word := range chirpWords {
		for _, cuss := range profaneWords {
			lowered := strings.ToLower(word)
			if lowered == cuss {
				chirpWords[i] = "****"
			}
		}
	}

	chirp.responseBody.CleanedBody = strings.Join(chirpWords, " ")
}
