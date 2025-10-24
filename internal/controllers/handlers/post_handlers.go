package handlers

import (
	"fmt"
	"net/http"
)

func TestReceiveFormHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("received request from form!")
}

func TestUserCreateHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("username: %s password: %s", email, password)))
}
