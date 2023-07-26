package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) add_user(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email string `json:"email"`
		Psw   string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	req := request{}
	err := decoder.Decode(&req)
	if err != nil {
		respondWithError(w, 400, "Something went wrong")
		return
	}

	responseBody, err := cfg.DB.CreateUser(req.Email, req.Psw)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	respondWithJSON(w, 201, responseBody)
}
