package main

import (
	"fmt"
	"net/http"
)

func main() {

	req, err := http.NewRequest("GET", "http://localhost:8000", nil)
	client := http.DefaultClient
	_, err = client.Do(req)
	fmt.Println("error returned: ", err)
}
