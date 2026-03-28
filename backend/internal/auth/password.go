package auth

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		log.Println("Failed to hash the password: ", err)
		return "", err
	}

	return string(hashBytes), nil

}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	if err != nil {
		log.Println("Failed to compare the password: ", err)
		return false
	}
	return true
}
