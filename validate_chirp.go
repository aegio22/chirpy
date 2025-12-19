package main

import (
	"encoding/json"
	"net/http"
)

type Chirp struct {
	Body string `json:"body"`
}

type ValidationErr struct {
	Error string `json:"error"`
}

type ChirpValidity struct {
	Valid bool `json:"valid"`
}

func handlerValidateChirp(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	chirp := Chirp{}
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

	validityResponse, err := json.Marshal(ChirpValidity{Valid: true})
	if err != nil {
		valErr, _ := json.Marshal(ValidationErr{Error: "Error marshaling validity response"})

		rw.WriteHeader(400)
		rw.Write(valErr)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(200)
	rw.Write(validityResponse)

}
