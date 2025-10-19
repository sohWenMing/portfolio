package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func main() {
	randBytes := make([]byte, 32)
	_, err := rand.Read(randBytes)
	if err != nil {
		fmt.Println("error occured: ", err)
	}
	encoded := base64.StdEncoding.EncodeToString(randBytes)
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		fmt.Println("error occured: ", err)
	}
	fmt.Println("encoded", encoded)
	fmt.Println("length of decoded: ", len(decoded))
}
