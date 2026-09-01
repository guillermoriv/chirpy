package main

import (
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/guillermoriv/chirpy/internal/database"
)

func databaseChirpToChirp(c database.Chirp) Chirp {
	return Chirp{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Body:      c.Body,
		UserID:    c.UserID,
	}
}

func databaseChirpsToChirps(chirps []database.Chirp, sortOrder string) []Chirp {
	result := make([]Chirp, len(chirps))
	for i, c := range chirps {
		result[i] = databaseChirpToChirp(c)
	}
	if sortOrder == "desc" {
		sort.Slice(result, func(i, j int) bool {
			return result[i].CreatedAt.After(result[j].CreatedAt)
		})
	}

	return result
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	authorID := uuid.Nil
	authorIDString := r.URL.Query().Get("author_id")
	sortOrder := r.URL.Query().Get("sort")

	if authorIDString != "" {
		var err error
		authorID, err = uuid.Parse(authorIDString)
		if err != nil {
			respondWithError(w, "invalid id for author", http.StatusBadRequest, err)
			return
		}
	}

	var (
		chirps []database.Chirp
		err    error
	)

	if authorID != uuid.Nil {
		chirps, err = cfg.db.GetChirpsByAuthorID(r.Context(), authorID)
	} else {
		chirps, err = cfg.db.GetChirps(r.Context())
	}

	if err != nil {
		respondWithError(w, "couldn't retrieve chirps", http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, databaseChirpsToChirps(chirps, sortOrder), http.StatusOK)
}
