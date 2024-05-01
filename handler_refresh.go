package main

import (
	"net/http"
	"time"

	"github.com/ross-brown/chirpy-go/internal/auth"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find token")
		return
	}

	user, err := cfg.DB.GetUserByRefreshToken(refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find user for refresh token")
		return
	}

	accessToken, err := auth.CreateToken(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
}
