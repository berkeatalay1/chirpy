package main

import (
	"encoding/json"
	"net/http"
)

type request struct {
	Body string `json:"body"`
}

type response struct {
	Valid bool `json:"valid"`
}

type error struct {
	Error string `json:"error"`
}

func validate_chirp(w http.ResponseWriter, r *http.Request) {
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
	respBody := response{
		Valid: true,
	}
	respondWithJSON(w, 200, respBody)
}
