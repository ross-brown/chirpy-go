package main

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/ross-brown/chirpy-go/internal/database"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	authorID := r.URL.Query().Get("author_id")
	sortParam := r.URL.Query().Get("sort")

	var dbChirps []database.Chirp
	var err error

	if authorID != "" {
		authorIDInt, err := strconv.Atoi(authorID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Couldn't convert string to int for author ID")
			return
		}

		dbChirps, err = cfg.DB.GetUserChirps(authorIDInt)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps")
			return
		}
	} else {
		dbChirps, err = cfg.DB.GetChirps()
	}

	if sortParam == "" {
		sortParam = "asc"
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps")
		return
	}

	sort.Slice(dbChirps, func(i, j int) bool {
		if sortParam == "desc" {
			return dbChirps[i].ID > dbChirps[j].ID
		}
		return dbChirps[i].ID < dbChirps[j].ID
	})

	respondWithJSON(w, http.StatusOK, dbChirps)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := strconv.Atoi(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't convert ID string to int")
		return
	}

	dbChirp, err := cfg.DB.GetChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp")
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:       dbChirp.ID,
		Body:     dbChirp.Body,
		AuthorID: dbChirp.AuthorID,
	})
}
