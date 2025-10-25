package passwordhashing

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var ErrPasswordTooShort = errors.New("password would be too short to be a valid password")

func GetErrorPasswordTooShort(password string) error {
	return fmt.Errorf("%s would be too short to be a valid error %w", password, ErrPasswordTooShort)
}

var ErrInvalidPassword = errors.New("password does not match the one listed in credentials")

func GetErrInvalidPassword(password string) error {
	return fmt.Errorf("%s does not match the one listen in credentials %w", password, ErrInvalidPassword)
}

var ErrPasswordTooLong = errors.New("password entered is too long")

func GetErrPasswordTooLong(password string) error {
	return fmt.Errorf("%s does not match the one listen ir credentials %w", password, ErrPasswordTooLong)
}

func HandlePasswordHashingError(password string, err error) error {
	if errors.Is(err, bcrypt.ErrHashTooShort) {
		return GetErrorPasswordTooShort(password)
	}
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return GetErrInvalidPassword(password)
	}
	if errors.Is(err, bcrypt.ErrPasswordTooLong) {
		return GetErrPasswordTooLong(password)
	}
	return err
}
