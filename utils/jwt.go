package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JwtKey = []byte("your_secret_key") // store this securely later (e.g., .env)

func GenerateJWT(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(), // expires in 24 hrs
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtKey)
}
