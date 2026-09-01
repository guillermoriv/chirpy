package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/guillermoriv/chirpy/internal/auth"
	"github.com/guillermoriv/chirpy/internal/database"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
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

	chirpIDString := r.PathValue("chirpID")

	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, "invalid uuid passed", http.StatusInternalServerError, err)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, "chirp not found", http.StatusNotFound, err)
		return
	}

	if chirp.UserID != userID {
		respondWithError(w, "couldn't delete chirp", http.StatusForbidden, err)
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), database.DeleteChirpParams{
		UserID: userID,
		ID:     chirpID,
	})
	if err != nil {
		respondWithError(w, "couldn't delete chirp", http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
