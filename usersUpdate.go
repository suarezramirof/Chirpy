package main

import (
	"encoding/json"
	"net/http"
	"github.com/suarezramirof/Chirpy/internal/auth"
	"strconv"
)

func (cfg *apiConfig) updateUser(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	
	userId, err := auth.CheckJWT(token, cfg.jwtSecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	hashedPassword, err := auth.HashedPassword(params.Password)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	numId, err := strconv.Atoi(userId)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := cfg.DB.UpdateUser(numId, params.Email, hashedPassword)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}