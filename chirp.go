package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func (cfg *apiConfig) chirpHandler(w http.ResponseWriter, r *http.Request) {
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

	cleanedBody := cleanBody(params.Body)

	chirp, err := cfg.DB.CreateChirp(cleanedBody)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create chirp")
		return
	}
	respondWithJSON(w, http.StatusCreated, chirp)
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

func (cfg *apiConfig) chirpsGetter(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get chirps")
		return
	}
	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) chirpGetter(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	numId, err := strconv.Atoi(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get chirp")
		return
	}
	// get chirp from db
	chirp, err := cfg.DB.GetChirp(numId)

	if err != nil {
		if err.Error() == "chirp not found" {
			respondWithError(w, http.StatusNotFound, "Chirp not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Could not get chirp")
		return
	}

	respondWithJSON(w, http.StatusOK, chirp)
}
