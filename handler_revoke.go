package main

import (
	"net/http"

	"github.com/ross-brown/chirpy-go/internal/auth"
)

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find refresh token")
		return
	}

	err = cfg.DB.RevokeRefreshToken(refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke refresh token")
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}
