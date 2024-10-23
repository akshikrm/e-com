package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

func CreateJwt(id int) (string, error) {
	claims := jwt.MapClaims{
		"exp": jwt.NewNumericDate(time.Unix(1516239022, 0)),
		"sub": id,
		"iat": time.Now(),
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}
