package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, "couldn't decode parameters", http.StatusInternalServerError, err)
		return
	}

	const maxChirpLength = 140
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	if len(params.Body) > maxChirpLength {
		respondWithError(w, "chirp too long", http.StatusBadRequest, nil)
		return
	}

	bodyResp := struct {
		CleanedBody string `json:"cleaned_body"`
	}{
		CleanedBody: getCleanedBody(params.Body, badWords),
	}

	respondWithJSON(w, bodyResp, http.StatusOK)
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}
