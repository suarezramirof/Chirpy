package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) upgradeUser(w http.ResponseWriter, r *http.Request) {

	type data struct {
		UserId int `json:"user_id"`
	}

	type body struct {
		Event string `json:"event"`
		Data  data   `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := body{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, 204, nil)
		return
	}

	userId := params.Data.UserId

	err = cfg.DB.UpgradeUser(userId)

	if err != nil {
		if err.Error() == "user not found" {
			respondWithJSON(w, 404, nil)
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, 204, nil)
}
