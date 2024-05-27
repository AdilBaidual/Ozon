package secure

import (
	"golang.org/x/crypto/bcrypt"
)

func CalculateHash(message string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(message), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CompareHash(hash, message string) bool {
	byteHash := []byte(hash)
	err := bcrypt.CompareHashAndPassword(byteHash, []byte(message))

	return err == nil
}
