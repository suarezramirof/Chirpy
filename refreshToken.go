package main

import (
	"net/http"

	"github.com/suarezramirof/Chirpy/internal/auth"
)

func (cfg *apiConfig) refreshToken(w http.ResponseWriter, r *http.Request) {
	refToken, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	userId, err := cfg.DB.CheckRefreshToken(refToken)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	jwt, err := auth.MakeJWT(userId, cfg.jwtSecret)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	type Token struct {
		Token string `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, Token{Token: jwt})
}
