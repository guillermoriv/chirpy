package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, "couldn't retrieve chirps", http.StatusInternalServerError, err)
		return
	}

	aggChirps := []Chirp{}

	for _, chirp := range chirps {
		aggChirps = append(aggChirps, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	respondWithJSON(w, aggChirps, http.StatusOK)
}
