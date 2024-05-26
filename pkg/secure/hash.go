package secure

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func CalculateSha512Signature(secret string, message string) string {
	mac := hmac.New(sha512.New, []byte(secret))
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

func CalculateMD5Signature(value, accessKey, accessID string) string {
	token := fmt.Sprintf("%s:%s:%s", value, accessKey, accessID)
	tokenHash := md5.Sum([]byte(token)) //nolint:gosec
	return hex.EncodeToString(tokenHash[:])
}

func CalculateSha256Signature(secret string, message string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	signature := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(signature)
}

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
