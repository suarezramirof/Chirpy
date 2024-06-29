package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"github.com/suarezramirof/Chirpy/internal/auth"
	"github.com/suarezramirof/Chirpy/shared"
)



func (cfg *apiConfig) chirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	userId, err := auth.CheckJWT(token, cfg.jwtSecret)

	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	cleanedBody := cleanBody(params.Body)

	chirp, err := cfg.DB.CreateChirp(cleanedBody, userIdInt)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create chirp")
		return
	}
	respondWithJSON(w, http.StatusCreated, shared.Chirp{
		Id: chirp.Id,
		Body: cleanedBody,
		AuthorId: userIdInt,
	})
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