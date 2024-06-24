package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func chirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	type success struct {
		Cleaned_body string `json:"cleaned_body"`
	}

	respondWithJSON(w, 200, success{cleanBody(params.Body)})
}

func cleanBody(body string) string {
	splitString := strings.Split(body, " ")
	for i, word := range splitString {
		lowerCase := strings.ToLower(word)
		if lowerCase == "kerfuffle" || lowerCase == "sharbert" || lowerCase == "fornax" {
			splitString[i] = "****"
		}
	}
	joinedString := strings.Join(splitString, " ")
	return joinedString
}
