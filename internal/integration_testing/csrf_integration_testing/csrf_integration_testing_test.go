package csrfintegrationtesting

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	loadenv "github.com/sohWenMing/portfolio/internal/env"
	csrf_protect "github.com/sohWenMing/portfolio/internal/security/csrf_protect"
)

var envGetter *loadenv.EnvGetter

// loads the envGetter into a global variable, to be used across all tests
func TestMain(m *testing.M) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	cleanedPath := cwd
	for range 3 {
		cleanedPath = filepath.Dir(cleanedPath)
	}
	envPath := fmt.Sprintf("%s/.env", cleanedPath)
	envGetter, err = loadenv.LoadEnv(envPath)
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestGetCSRFKeyFromEnv(t *testing.T) {
	expectedLengthBytes := 32
	csrfKeyStruct := envGetter.GetCSRFKey()
	if len(csrfKeyStruct.KeyBytes) != expectedLengthBytes {
		t.Errorf("want :%d\ngot%d", expectedLengthBytes, len(csrfKeyStruct.KeyBytes))
	}
}

func TestGetCSRFKeyFromCSRFProtect(t *testing.T) {
	expectedLengthBytes := 32
	csrfKey := csrf_protect.GetCSRFKey(envGetter)
	if len(csrfKey) != expectedLengthBytes {
		t.Errorf("want :%d\ngot%d", expectedLengthBytes, len(csrfKey))
	}
}
