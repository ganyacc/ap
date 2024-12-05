package util

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JwtClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// we can store secret in env or any configuration file. Currently we are storing in application itself for just simplicity.
var (
	SecretKey        = []byte("zxgasd4598d6as4das4d5sa4d4")
	RefreshSecretKey = []byte("gsdfs5sd4f4s5d4f54s5df")
)

func GenerateJwtToken(email string) (string, error) {

	claims := JwtClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "go-application",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token expiry after 24 hours
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token and return it as a string
	return token.SignedString(SecretKey)
}

func GenerateRefreshToken(email string) (string, error) {
	claims := &JwtClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "go-application",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * 24 * time.Hour)), // Refresh token expires in 10 days
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(RefreshSecretKey)
}
