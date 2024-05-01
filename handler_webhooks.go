package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ross-brown/chirpy-go/internal/auth"
)

func (cfg *apiConfig) handlerWebhooks(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		}
	}

	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil || apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Invalid or missing API Key")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode request body")
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusOK, struct{}{})
		return
	}

	_, err = cfg.DB.UpgradeChirpyRed(params.Data.UserID)
	if err != nil {
		if errors.Is(err, errors.New("resource does not exist")) {
			respondWithError(w, http.StatusNotFound, "Couldn't find user")
			return
		}

		respondWithError(w, http.StatusInternalServerError, "Couldn't upgrade user")
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}
