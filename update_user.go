package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/berkeatalay/chirpy/internal/auth"
)

func (cfg *apiConfig) update_user(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email string `json:"email"`
		Psw   string `json:"password"`
	}
	type response struct {
		Email  string `json:"email"`
		UserId int32  `json:"id"`
	}
	decoder := json.NewDecoder(r.Body)
	req := request{}
	err := decoder.Decode(&req)
	if err != nil {
		respondWithError(w, 400, "Something went wrong")
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 400, "Error Getting Header")
		return
	}

	userId, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}
	userIDInt, err := strconv.Atoi(userId)
	if err != nil {
		respondWithError(w, 400, "Error Parsing Int")
		return
	}
	psw, err := auth.HashPassword(req.Psw)
	if err != nil {
		respondWithError(w, 400, "Something went wrong")
		return
	}
	user, err := cfg.DB.UpdateUser(userIDInt, req.Email, psw)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	respondWithJSON(w, 200, response{
		Email:  user.Email,
		UserId: int32(userIDInt),
	})
}
