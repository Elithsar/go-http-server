package main

import (
	"encoding/json"
	"net/http"
)

func ChirpIn(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Body string `json:"body"`
	}
	type response struct {
		Valid bool `json:"valid,omitempty"`
	}
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	var message request
	err := decoder.Decode(&message)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	MAX_SIZE := 140
	if len(message.Body) > MAX_SIZE {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Valid: true,
	})
}
