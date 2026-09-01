package main

import (
	"net/http"
	"time"

	"github.com/guillermoriv/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, "not token available in the headers", http.StatusBadRequest, err)
		return
	}

	refreshToken, err := cfg.db.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, "unauthorized", http.StatusUnauthorized, err)
		return
	}

	if time.Now().UTC().After(refreshToken.ExpiresAt) || refreshToken.RevokedAt.Valid {
		respondWithError(w, "unauthorized", http.StatusUnauthorized, nil)
		return
	}

	accessToken, err := auth.MakeJWT(refreshToken.UserID, cfg.secret, 1*time.Hour)
	if err != nil {
		respondWithError(w, "couldn't generate the access token", http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, response{Token: accessToken}, http.StatusOK)
}
