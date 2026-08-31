package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/guillermoriv/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds *int   `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, "couldn't decode parameters", http.StatusBadRequest, err)
		return
	}

	expirationTime := 1 * time.Hour
	if params.ExpiresInSeconds != nil {
		requestedDuration := time.Duration(*params.ExpiresInSeconds) * time.Second
		if requestedDuration < expirationTime {
			expirationTime = requestedDuration
		}
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, "incorrect email or password", http.StatusUnauthorized, err)
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !match {
		respondWithError(w, "incorrect email or password", http.StatusUnauthorized, err)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, expirationTime)
	if err != nil {
		respondWithError(w, "couldn't generate token for user", http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     token,
	}, http.StatusOK)
}
