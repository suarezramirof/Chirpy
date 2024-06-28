package main

import (
	"net/http"

	"github.com/suarezramirof/Chirpy/internal/auth"
)

func (cfg *apiConfig) revokeToken(w http.ResponseWriter, r *http.Request) {
	refToken, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	_, err = cfg.DB.CheckRefreshToken(refToken)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	err = cfg.DB.DeleteRefreshToken(refToken)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, 204, "")
}
