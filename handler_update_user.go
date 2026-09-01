package main

import (
	"encoding/json"
	"net/http"

	"github.com/guillermoriv/chirpy/internal/auth"
	"github.com/guillermoriv/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, "couldn't decode parameters", http.StatusInternalServerError, err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, "couldn't find authorization", http.StatusUnauthorized, err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, "couldn't validate the JWT", http.StatusUnauthorized, err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, "couldn't hash this password", http.StatusInternalServerError, err)
		return
	}

	user, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
		ID:             userID,
	})
	if err != nil {
		respondWithError(w, "couldn't update user information", http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.CreatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}, http.StatusOK)
}
