package jwtToken

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("your_secret_key")

// GenerateToken creates a new JWT token
func GenerateToken(userName string) (string, error) {
	tokenClaims := jwt.MapClaims{
		"userName": userName,
		"exp":      time.Now().Add(time.Hour * 1).Unix(), // Token expiration
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	return token.SignedString(secretKey)
}

var ErrInvalidToken = errors.New("invalid token")
