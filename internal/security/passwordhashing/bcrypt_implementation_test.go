package passwordhashing

import "testing"

func TestBCryptImplementation(t *testing.T) {
	type test struct {
		name                 string
		password             string
		attemptCheckPassword string
		isErrExpected        bool
	}
	tests := []test{
		{
			"test base case, should pass",
			"password123",
			"password123",
			false,
		},
		{
			"test failure case",
			"password123",
			"password234",
			true,
		},
	}

	var bcrypt = &BcryptInterfacer{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hashedPassword, err := bcrypt.HashPassword(test.password)
			if err != nil {
				t.Errorf("didn't expect errror got %v", err)
				return
			}
			err = bcrypt.CheckPassword(test.attemptCheckPassword, hashedPassword)
			switch test.isErrExpected {
			case true:
				if err == nil {
					t.Errorf("expected error, didn't get one")
					return
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
