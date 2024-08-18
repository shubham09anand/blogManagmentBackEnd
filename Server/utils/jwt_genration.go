package jwtToken

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

// Claims defines the structure of JWT claims
type Claims struct {
	Plaintext string `json:"plaintext"`
	jwt.RegisteredClaims
}

// Encrypt generates a JWT token with the given plaintext
func Encrypt(key []byte, plaintext string) (string, error) {
	// Create the claims with plaintext and expiration time
	claims := &Claims{
		Plaintext: plaintext,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Second)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
