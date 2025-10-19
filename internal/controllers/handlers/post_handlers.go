package handlers

import (
	"fmt"
	"net/http"
)

func TestReceiveFormHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("received request from form!")
}
