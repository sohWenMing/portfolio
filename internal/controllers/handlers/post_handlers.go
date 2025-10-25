package handlers

import (
	"fmt"
	"net/http"

	businesslogic "github.com/sohWenMing/portfolio/internal/business_logic"
)

func TestReceiveFormHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("received request from form!")
}

func TestUserCreateHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	w.WriteHeader(http.StatusOK)
	w.Write(fmt.Appendf(nil, "username: %s password: %s", email, password))
}

func CreateUserHandler(us *businesslogic.UserService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.FormValue("email")
		password := r.FormValue("passwrod")
		fmt.Println("email: ", email)
		fmt.Println("password: ", password)
		//TODO: finish handler

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))

	}

}
