package csrfintegrationtesting

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	loadenv "github.com/sohWenMing/portfolio/internal/env"
)

var envGetter *loadenv.EnvGetter

func TestMain(m *testing.M) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("err occured: ", err)
	}
	fmt.Println("cwd: ", cwd)
	cleanedPath := cwd
	for range 3 {
		cleanedPath = filepath.Dir(cleanedPath)
	}
	envPath := fmt.Sprintf("%s/.env", cleanedPath)
	envGetter, err = loadenv.LoadEnv(envPath)
	fmt.Println("cleanedPath: ", cleanedPath)

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestGetCSRFKeyFromEnv(t *testing.T) {
	csrfKeyStruct := envGetter.GetCSRFKey()
	fmt.Println("keyString: ", csrfKeyStruct.KeyString)
	fmt.Println("Length of keyBytes: ", len(csrfKeyStruct.KeyBytes))
}
