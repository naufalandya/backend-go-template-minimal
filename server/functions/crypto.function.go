package functions

import (
	"encoding/base64"

	"golang.org/x/crypto/argon2"
)

func HashPassword(password string) (string, error) {
	salt := []byte("cutie_salt")
	hash := argon2.Key([]byte(password), salt, 1, 64*1024, 4, 32)
	hashedPassword := base64.RawStdEncoding.EncodeToString(hash)
	return hashedPassword, nil
}
