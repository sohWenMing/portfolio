package mainservertesting

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"

	dbinfra "github.com/sohWenMing/portfolio/internal/db_infra"
	dbinterface "github.com/sohWenMing/portfolio/internal/db_interface"
	integration "github.com/sohWenMing/portfolio/internal/integration"
	"github.com/sohWenMing/portfolio/internal/security/passwordhashing"
)

var appDb *dbinfra.AppDB
var appServices *integration.Services
var server *http.Server
var addr = ":8000"

func TestMain(m *testing.M) {
	returnedAppDB, handler, returnedServices, err := integration.Run(true, "../../../.env")
	if err != nil {
		panic(err)
	}
	appDb = returnedAppDB
	appServices = returnedServices
	server = &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("server could not start")
			panic(err)
		}
		fmt.Println("server started listening on port: ", server.Addr)
	}()
	testPing()

	exitCode := m.Run()
	ctx, cancel := context.WithTimeout(context.Background(), 5&time.Second)
	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("error during server shutdown: ", err)
		panic(err)
	}
	cancel()
	err = appDb.DB.Close()
	if err != nil {
		panic(err)
	}
	os.Exit(exitCode)
}

func testPing() {
	curCount := 1
	client := http.DefaultClient
	for i := curCount; i <= 5; i++ {
		fmt.Println("current count in testPing: ", curCount)
		err := ping(client)
		if err == nil {
			return
		}
		time.Sleep(1 * time.Second)
		curCount += 1
	}
}

func ping(client *http.Client) error {
	req, err := http.NewRequest("GET", "http://localhost:8000", nil)
	if err != nil {
		panic(err)
	}
	res, err := client.Do(req)
	if err == nil && res.StatusCode == 200 {
		fmt.Println("Setup was ok: ping to server got status code 200")
		return nil
	}
	return err
}

func TestCreateUser(t *testing.T) {
	id, err := appServices.UserService.CreateUser(
		"wenming.soh@gmail.com",
		"Holoq123holoq123")
	if err != nil {
		t.Errorf("didn't expect error, got %v", err)
		return
	}
	numrows, err := appServices.UserService.DeleteUserById(id)
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
		id, err := appServices.UserService.CreateUser(test.email, test.createdPassword)
		if err != nil {
			t.Errorf("didn't expect error, got %v", err)
			return
		}
		createdId := id
		details, err := appServices.UserService.GetUserDetailsById(createdId)
		if err != nil {
			t.Errorf("didn't expect error, got %v", err)
			return
		}
		err = appServices.UserService.CheckPassword(
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
		numrows, err := appServices.UserService.DeleteUserById(id)
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
	details, err := appServices.UserService.GetUserDetailsById(0)

	if !errors.Is(err, expectedErr) {
		t.Errorf("expected %v\n got %v", err, expectedErr)
	}
	if !reflect.DeepEqual(details, dbinterface.UserDetails{}) {
		t.Errorf("got %v\n want %v", details, dbinterface.UserDetails{})
	}
}

func TestDeleteNonExistentUser(t *testing.T) {

	numrows, err := appServices.UserService.DeleteUserById(0)
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
