package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Chirp struct {
	Body         string `json:"body"`
	responseBody struct {
		CleanedBody string `json:"cleaned_body"`
	}
}

type ValidationErr struct {
	Error string `json:"error"`
}

type ChirpValidity struct {
	Valid bool `json:"valid"`
}

func handlerValidateChirp(rw http.ResponseWriter, req *http.Request) {

	chirp := Chirp{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&chirp)
	if err != nil {
		valErr, _ := json.Marshal(ValidationErr{
			Error: "Error decoding chirp body",
		})

		rw.WriteHeader(400)
		rw.Write(valErr)
		return
	}

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
	chirpJson, err := json.Marshal(chirp.responseBody)
	if err != nil {
		valErr, _ := json.Marshal(ValidationErr{
			Error: "Error marshalling cleaned body",
		})

		rw.WriteHeader(400)
		rw.Write(valErr)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(200)
	rw.Write(chirpJson)

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
