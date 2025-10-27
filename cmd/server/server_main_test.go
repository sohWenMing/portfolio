package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"

	dbinfra "github.com/sohWenMing/portfolio/internal/db_infra"
	dbinterface "github.com/sohWenMing/portfolio/internal/db_interface"
	"github.com/sohWenMing/portfolio/internal/security/passwordhashing"
)

var appDb *dbinfra.AppDB
var appServices *services

func TestMain(m *testing.M) {
	returnedAppDB, _, returnedServices, err := Run(true, "../../.env")
	if err != nil {
		panic(err)
	}
	appDb = returnedAppDB
	appServices = returnedServices
	exitCode := m.Run()
	err = appDb.DB.Close()
	if err != nil {
		panic(err)
	}
	os.Exit(exitCode)
}

func TestCreateUser(t *testing.T) {
	id, err := appServices.userservice.CreateUser(
		"wenming.soh@gmail.com",
		"Holoq123holoq123")
	if err != nil {
		t.Errorf("didn't expect error, got %v", err)
		return
	}
	numrows, err := appServices.userservice.DeleteUserById(id)
	if err != nil {
		t.Errorf("didn't expect error, got %v", err)
		unwrapAndPrintError(err)
		return
	}
	numRowsExpectDeleted := 1
	if numrows != int64(numRowsExpectDeleted) {
		t.Errorf("got %d\nwant %d", numrows, numRowsExpectDeleted)
		return
	}
}

func TestGetUserDetails(t *testing.T) {
	type test struct {
		testName          string
		email             string
		createdPassword   string
		attemptedPassword string
		isErrExpected     bool
		expectedError     error
	}

	tests := []test{
		{
			"basic test, should pass",
			"wenming.soh@gmail.com",
			"Holoq123holoq123",
			"Holoq123holoq123",
			false,
			nil,
		},
		{
			"basic test, should pass",
			"wenming.soh@gmail.com",
			"Holoq123holoq123",
			"fail",
			true,
			passwordhashing.ErrInvalidPassword,
		},
	}

	for _, test := range tests {
		id, err := appServices.userservice.CreateUser(test.email, test.createdPassword)
		if err != nil {
			t.Errorf("didn't expect error, got %v", err)
			return
		}
		createdId := id
		details, err := appServices.userservice.GetUserDetailsById(createdId)
		if err != nil {
			t.Errorf("didn't expect error, got %v", err)
			return
		}
		err = appServices.userservice.CheckPassword(
			test.attemptedPassword,
			details.HashedPassword,
		)
		switch test.isErrExpected {
		case true:
			verifyErr := verifyError(err, test.expectedError)
			if verifyErr != nil {
				t.Errorf("didn't expect error on verify error, got %v", err)
				return
			}
		case false:
			if err != nil {
				t.Errorf("didn't expect error, got %v", err)
				return
			}
		}
		numrows, err := appServices.userservice.DeleteUserById(id)
		if err != nil {
			t.Errorf("didn't expect error, got %v", err)
			return
		}
		numRowsExpectDeleted := 1
		if numrows != int64(numRowsExpectDeleted) {
			t.Errorf("got %d\nwant %d", numrows, numRowsExpectDeleted)
			return
		}
	}
}

func TestGetNonExistentUser(t *testing.T) {
	expectedErr := sql.ErrNoRows
	details, err := appServices.userservice.GetUserDetailsById(0)

	if !errors.Is(err, expectedErr) {
		t.Errorf("expected %v\n got %v", err, expectedErr)
	}
	if !reflect.DeepEqual(details, dbinterface.UserDetails{}) {
		t.Errorf("got %v\n want %v", details, dbinterface.UserDetails{})
	}
}

func TestDeleteNonExistentUser(t *testing.T) {

	numrows, err := appServices.userservice.DeleteUserById(0)
	if err != nil {
		t.Errorf("didn't expect error, got %v", err)
		unwrapAndPrintError(err)
		return
	}
	numRowsExpectDeleted := 0
	if numrows != int64(numRowsExpectDeleted) {
		t.Errorf("got %d\nwant %d", numrows, numRowsExpectDeleted)
		return
	}
}

func verifyError(input error, expected error) error {
	errorToUnwrap := input
	for {
		if errors.Is(errorToUnwrap, expected) {
			return nil
		}
		if errors.Unwrap(errorToUnwrap) == nil {
			return fmt.Errorf(
				"input error is not and does not wrap expected error\ngot %v want %v", input, expected)
		}
		errorToUnwrap = errors.Unwrap(errorToUnwrap)
	}
}

func unwrapAndPrintError(err error) {
	currErr := err
	layer := 1
	for {
		fmt.Println("Error Layer: ", layer)
		fmt.Printf("Error Type: %T\n", currErr)
		fmt.Println("error: ", currErr)
		if errors.Unwrap(currErr) == nil {
			return
		} else {
			layer += 1
			currErr = errors.Unwrap(err)
		}
	}
}
