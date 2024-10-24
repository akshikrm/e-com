package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateJwt(id int) (string, error) {
	claims := jwt.MapClaims{
		"exp": jwt.NewNumericDate(time.Unix(1516239022, 0)),
		"sub": id,
		"iat": time.Now(),
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwt, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to generate token")
	}
	return jwt, nil
}

func ValidateJWT(token string) (*jwt.Token, error) {
	return jwt.Parse(token,
		func(token *jwt.Token) (interface{}, error) {
			secret := os.Getenv("JWT_SECRET")
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		},
	)
}
