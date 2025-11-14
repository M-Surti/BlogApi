package utils

import (


	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	HashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	

	return string(HashPassword), nil
}

func ComparePassword(plainPass, Hashpass string) error {
	err := bcrypt.CompareHashAndPassword([]byte(Hashpass), []byte(plainPass))
	if err != nil {
		return err
	}
	return nil
}