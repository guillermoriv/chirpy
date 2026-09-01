package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/guillermoriv/chirpy/internal/auth"
	"github.com/guillermoriv/chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, "couldn't decode parameters", http.StatusBadRequest, err)
		return
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

	accessToken, err := auth.MakeJWT(user.ID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, "couldn't generate token for user", http.StatusInternalServerError, err)
		return
	}

	refreshToken := auth.MakeRefreshToken()

	err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{Token: refreshToken, UserID: user.ID, ExpiresAt: time.Now().UTC().AddDate(0, 0, 60)})
	if err != nil {
		respondWithError(w, "couldn't save token for user", http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, response{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		IsChirpyRed:  user.IsChirpyRed,
		Token:        accessToken,
		RefreshToken: refreshToken,
	}, http.StatusOK)
}
