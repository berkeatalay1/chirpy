package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/berkeatalay/chirpy/internal/auth"
	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) login(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email      string `json:"email"`
		Psw        string `json:"password"`
		ExpSeconds int    `json:"expires_in_seconds"`
	}
	defaultExpiration := 60 * 60 * 24
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

	if req.ExpSeconds == 0 {
		req.ExpSeconds = defaultExpiration
	} else if req.ExpSeconds > defaultExpiration {
		req.ExpSeconds = defaultExpiration
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Duration(req.ExpSeconds) * time.Second)),
		Subject:   fmt.Sprintf("%d", user.Id),
	})

	key, err := base64.StdEncoding.DecodeString(cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	signedToken, err := token.SignedString(key)
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
	rsp.Token = signedToken

	respondWithJSON(w, 200, rsp)
}
