package auth

import (
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrNoAuthHeaderIncluded = errors.New("not auth header included in request")

// HashPassword -
func HashPassword(password string) (string, error) {
	dat, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(dat), nil
}

// CheckPasswordHash -
func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		return "", errors.New("malformed authorization header")
	}
	return splitAuth[1], nil
}

func ValidateJWT(tokenString, tokenSecret string) (string, error) {
	claimsStruct := jwt.RegisteredClaims{}
	tokenSecretString, err := base64.StdEncoding.DecodeString(tokenSecret)
	if err != nil {
		log.Printf("0 %s", err)
		return "", err
	}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecretString), nil },
	)
	if err != nil {
		log.Printf("1 %s", err)
		return "", err
	}
	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		log.Printf("2 %s", err)
		return "", err
	}
	expiresAt, err := token.Claims.GetExpirationTime()
	if err != nil {
		log.Printf("3 %s", err)
		return "", err
	}
	if expiresAt.Before(time.Now().UTC()) {
		log.Printf("4 %s", err)
		return "", errors.New("JWT is expired")
	}
	return userIDString, nil
}
