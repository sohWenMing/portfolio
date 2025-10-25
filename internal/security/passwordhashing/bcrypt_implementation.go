package passwordhashing

import (
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

type BcryptInterfacer struct{}

// Hashes password, and returns base64 encoded string of hash
func (b *BcryptInterfacer) HashPassword(password string) (hash string, err error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	base64Encoding := base64.RawStdEncoding.EncodeToString(hashBytes)
	return base64Encoding, nil
}
func (b *BcryptInterfacer) CheckPassword(attemptedPassword string, hashedPassword string) (err error) {
	decodedPassword, err := base64.RawStdEncoding.DecodeString(hashedPassword)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword(decodedPassword, []byte(attemptedPassword))
	if err != nil {
		return err
	}
	return nil
}
