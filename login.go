package main

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/berkeatalay/chirpy/internal/auth"
	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) login(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email      string `json:"email"`
		Psw        string `json:"password"`
		ExpSeconds string `json:"expires_in_seconds"`
	}
	decoder := json.NewDecoder(r.Body)
	req := request{}
	err := decoder.Decode(&req)
	if err != nil {
		respondWithError(w, 400, "Something went wrong")
		return
	}

	user, isFound := cfg.DB.GetUserWithEmail(req.Email)
	if !isFound {
		respondWithError(w, 404, "Not Found User")
		return
	}
	checkPsw := auth.CheckPasswordHash(req.Psw, user.Password)
	if checkPsw != nil {
		respondWithError(w, 401, checkPsw.Error())
		return
	}

	expSeconds, err := strconv.Atoi(req.ExpSeconds)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:  "chirpy",
		Subject: string(user.Id),
		ExpiresAt: &jwt.NumericDate{
			Time: time.Now().UTC().Add(time.Duration(expSeconds) * time.Second),
		},
		IssuedAt: &jwt.NumericDate{
			Time: time.Now().UTC(),
		},
	})

	key, err := base64.StdEncoding.DecodeString(cfg.JwtSecret)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	token, err := jwt.SignedString(key)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	type response struct {
		Email string `json:"email"`
		Id    int32  `json:"id"`
		Token string `json:"token"`
	}
	rsp := response{}
	rsp.Email = user.Email
	rsp.Id = user.Id
	rsp.Token = token

	respondWithJSON(w, 200, rsp)
}
