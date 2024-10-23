package utils

import (
	"akshidas/e-com/pkg/model"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

func CreateJwt(u *model.User) (string, error) {
	claims := jwt.MapClaims{
		"exp": jwt.NewNumericDate(time.Unix(1516239022, 0)),
		"sub": u.ID,
		"iat": time.Now(),
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}
