package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password []byte) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func ValidateHash(hashedpassword []byte, plainTextPassword string) error {
	if err := bcrypt.CompareHashAndPassword(
		hashedpassword,
		[]byte(plainTextPassword)); err != nil {
		return Unauthorized
	}

	return nil
}
