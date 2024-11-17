package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

func CreateJwt(id int, role string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"exp":  now.Add(time.Hour * 24).Unix(),
		"sub":  id,
		"iat":  now.Unix(),
		"role": role,
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
