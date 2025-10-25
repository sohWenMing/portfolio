package passwordhashing

import (
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

type BcryptInterfacer struct{}

// Hashes password, and returns base64 encoded string of hash
func (b *BcryptInterfacer) HashPassword(password string) (hash string, err error) {
	if len(password) < 8 {
		return "", GetErrorPasswordTooShort(password)
	}
	if len(password) > 32 {
		return "", GetErrPasswordTooLong(password)
	}
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", HandlePasswordHashingError(password, err)
	}
	base64Encoding := base64.RawStdEncoding.EncodeToString(hashBytes)
	return base64Encoding, nil
}
func (b *BcryptInterfacer) CheckPassword(attemptedPassword string, hashedPassword string) (err error) {
	decodedPassword, err := base64.StdEncoding.DecodeString(hashedPassword)
	if err != nil {
		return HandlePasswordHashingError(attemptedPassword, err)
	}
	err = bcrypt.CompareHashAndPassword(decodedPassword, []byte(attemptedPassword))
	if err != nil {
		return HandlePasswordHashingError(attemptedPassword, err)
	}
	return nil
}
