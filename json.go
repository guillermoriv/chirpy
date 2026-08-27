package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, msg string, statusCode int, err error) {
	if err != nil {
		log.Println(err)
	}

	errorResp := struct {
		Error string `json:"error"`
	}{
		Error: msg,
	}

	respondWithJSON(w, errorResp, statusCode)
}

func respondWithJSON(w http.ResponseWriter, payload any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
	w.Write(dat)
}
