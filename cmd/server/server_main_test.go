package main

import (
	"fmt"
	"os"
	"testing"

	dbinfra "github.com/sohWenMing/portfolio/internal/db_infra"
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
	fmt.Println("appservices: ", appServices)
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
	err = appServices.userservice.DeleteUserById(id)
	if err != nil {
		t.Errorf("didn't expect error, got %v", err)
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
	}

	tests := []test{
		{
			"basic test, should pass",
			"wenming.soh@gmail.com",
			"Holoq123holoq123",
			"Holoq123holoq123",
			false,
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

	}
}
