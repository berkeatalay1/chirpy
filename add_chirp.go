package main

import (
	"encoding/json"
	"net/http"
)

type request struct {
	Body string `json:"body"`
}

type error struct {
	Error string `json:"error"`
}

func (cfg *apiConfig) add_chirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	req := request{}
	err := decoder.Decode(&req)
	if err != nil {
		respondWithError(w, 400, "Something went wrong")
		return
	}

	if len(req.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	res := replace_bad_words(req.Body)

	responseBody, err := cfg.DB.CreateChirp(res)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	respondWithJSON(w, 201, responseBody)
}
