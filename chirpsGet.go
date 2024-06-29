package main

import (
	"github.com/suarezramirof/Chirpy/internal/auth"
	"github.com/suarezramirof/Chirpy/shared"
	"net/http"
	"strconv"
	s "sort"
)

func (cfg *apiConfig) chirpsGetter(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.DB.GetChirps()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get chirps")
		return
	}
	
	authorId := r.URL.Query().Get("author_id")

	filteredChirps := []shared.Chirp{}

	if authorId != "" {
		numAuthorId, err := strconv.Atoi(authorId)
		if err != nil {
			respondWithError(w, 500, "Something went wrong")
			return
		}
		for _, chirp := range chirps {
			if chirp.AuthorId == numAuthorId {
				filteredChirps = append(filteredChirps, shared.Chirp{Id: chirp.Id, Body: chirp.Body, AuthorId: chirp.AuthorId})
			}
		}
	} else {
			filteredChirps = chirps
	}

	sort := r.URL.Query().Get("sort")
	if sort != "desc" {
		sort = "asc"
	}

	if sort == "desc" {
		s.Slice(filteredChirps, func(i, j int) bool {
			return filteredChirps[i].Id > filteredChirps[j].Id
		})
	} else {
		s.Slice(filteredChirps, func(i, j int) bool {
			return filteredChirps[i].Id < filteredChirps[j].Id
		})
	}
	
	respondWithJSON(w, http.StatusOK, filteredChirps)
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

func (cfg *apiConfig) deleteChirp(w http.ResponseWriter, r *http.Request) {
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

	numUserId, err := strconv.Atoi(userId)

	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	id := r.PathValue("id")
	numId, err := strconv.Atoi(id)

	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	chirps, err := cfg.DB.GetChirps()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get chirps")
		return
	}

	var chirpAuthor int

	for _, chirp := range chirps {
		if chirp.Id == numId {
			chirpAuthor = chirp.AuthorId
		}
	}

	if chirpAuthor == 0 {
		respondWithError(w, 404, "Chirp not found")
		return
	}

	if chirpAuthor != numUserId {
		respondWithError(w, 403, "Unauthorized")
		return
	}

	err = cfg.DB.DeleteChirp(numId)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not delete chirp")
		return
	}

	respondWithJSON(w, 204, "Chirp deleted")
}
