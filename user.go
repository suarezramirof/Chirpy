package main

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), 10)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// add user to database
	user, err := cfg.DB.AddUser(params.Email, hashedPassword)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, user)
}

func (cfg *apiConfig) loginUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type UserResponse struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
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
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(params.Password))

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	respondWithJSON(w, http.StatusOK, UserResponse{Id: user.Id, Email: user.Email})
}
