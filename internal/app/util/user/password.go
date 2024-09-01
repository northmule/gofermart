package user

import (
	"crypto/sha256"
	"fmt"
)

// PasswordHash хэш пароля
func PasswordHash(password string) string {
	hashAlg := sha256.New()
	hashAlg.Write([]byte(password))
	return fmt.Sprintf("%x", hashAlg.Sum(nil))
}
