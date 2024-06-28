package main

import (
	"encoding/json"
	"net/http"
	"github.com/suarezramirof/Chirpy/internal/auth"
)

func (cfg *apiConfig) loginUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds,omitempty"`
	}

	type UserResponse struct {
		Id           int    `json:"id"`
		Email        string `json:"email"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := cfg.DB.GetUser(params.Email)

	if err != nil {
		if err.Error() == "user not found" {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = auth.CompareHashAndPassword(user.Password, params.Password)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	tok, err := auth.MakeJWT(user.Id, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = cfg.DB.SaveRefreshToken(user.Id, refreshToken)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK,
		UserResponse{
			Id:           user.Id,
			Email:        user.Email,
			Token:        tok,
			RefreshToken: refreshToken})
}
