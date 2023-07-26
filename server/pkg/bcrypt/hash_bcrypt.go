package bcrypt

import (
	"golang.org/x/crypto/bcrypt"
)

func GeneratePassword(password string) (string, error) {

	hashByte, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}

	return string(hashByte), err
}

func ComparePassword(password, hashPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))

	return err == nil
}
