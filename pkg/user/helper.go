package user

import (
	"akshidas/e-com/pkg/types"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func createJwt(u *types.User) (string, error) {
	claims := jwt.MapClaims{
		"exp": jwt.NewNumericDate(time.Unix(1516239022, 0)),
		"sub": u.ID,
		"iat": time.Now(),
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func hashPassword(password []byte) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func validateHash(hashedpassword []byte, plainTextPassword string) error {

	if err := bcrypt.CompareHashAndPassword(
		hashedpassword,
		[]byte(plainTextPassword)); err != nil {
		return Unauthorized
	}

	return nil
}

func scanIntoUser(rows *sql.Rows) (*types.User, error) {
	user := &types.User{}
	err := rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		log.Printf("scan into user: %s", err)
	}

	return user, err
}
