package passwordhashing

import (
	"errors"
	"testing"
)

func TestBCryptImplementation(t *testing.T) {
	type test struct {
		name                  string
		password              string
		attemptCheckPassword  string
		isPWLengthErrExpected bool
		isFinalErrExpected    bool
		expectedErr           error
	}
	tests := []test{
		{
			"test base case, should pass",
			"password123",
			"password123",
			false,
			false,
			nil,
		},
		{
			"test failure case",
			"password123",
			"password234",
			false,
			true,
			ErrInvalidPassword,
		},
		{
			"test failure case password too short",
			"",
			"password234",
			true,
			false,
			ErrPasswordTooShort,
		},
		{
			"test failure case password too long",
			"1234567890|1234567890|1234567890|1234567890",
			"1234567890|1234567890|1234567890|1234567890",
			true,
			false,
			ErrPasswordTooLong,
		},
	}

	var bcrypt = &BcryptInterfacer{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hashedPassword, err := bcrypt.HashPassword(test.password)
			switch test.isPWLengthErrExpected {
			case true:
				isErrMatch, finalErr := checkErrUnwrapper(test.expectedErr, err)
				if !isErrMatch {
					t.Errorf("expected error did not match. Expected error: %v Final Error after unwrapping: %v",
						test.expectedErr, finalErr)
					return
				}
				return
			case false:
				if err != nil {
					t.Errorf("didn't expect error got %v", err)
					return
				}

			}
			err = bcrypt.CheckPassword(test.attemptCheckPassword, hashedPassword)
			switch test.isFinalErrExpected {
			case true:
				if err == nil {
					t.Errorf("expected error, didn't get one")
					return
				} else {
					isErrMatch, finalErr := checkErrUnwrapper(test.expectedErr, err)
					if !isErrMatch {
						t.Errorf("expected error did not match. Expected error: %v Final Error after unwrapping: %v",
							test.expectedErr, finalErr)
						return
					}
				}

			case false:
				if err != nil {
					t.Errorf("didn't expect errror got %v", err)
					return
				}

			}
		})
	}
}

func checkErrUnwrapper(expectedErr error, inputError error) (isExpectedErrorFound bool, finalError error) {
	isExpectedErrorFound = false
	errToUnwrap := inputError
	for {
		if errors.Is(errToUnwrap, expectedErr) {
			return true, nil
		} else if errors.Unwrap(errToUnwrap) == nil {
			return false, errToUnwrap
		} else {
			errToUnwrap = errors.Unwrap(errToUnwrap)
		}
	}
}
