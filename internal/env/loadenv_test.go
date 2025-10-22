package loadenv

import (
	"os"
	"reflect"
	"testing"
)

var envGetter *EnvGetter

func TestMain(m *testing.M) {
	loadedGetter, err := LoadEnv("../../.env")
	if err != nil {
		panic(err)
	}
	envGetter = loadedGetter
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestGetCSRFKey(t *testing.T) {
	keyData := envGetter.GetCSRFKey()
	checkCSRFKey := os.Getenv("CSRFKEY")
	if keyData.KeyString != checkCSRFKey {
		t.Errorf("got: %s\nwant %s", keyData.KeyString, checkCSRFKey)
		return
	}
	if len(keyData.KeyBytes) != 32 {
		t.Errorf("length of keybytes was not equal to 32")
		return
	}
}

func TestGetCSRFTrustedOrigins(t *testing.T) {
	expected := []string{
		"localhost:3000",
		"128.199.244.100",
	}
	csrfTrustedOrigins := envGetter.GetTrustedOrigins()
	if !reflect.DeepEqual(csrfTrustedOrigins, expected) {
		t.Errorf("got %v\nwant %v", csrfTrustedOrigins, expected)
	}
}

func TestGetDBString(t *testing.T) {
	expected := "host=localhost dbname=portfoliodb user=nindgabeet password=Holoq123holoq123 sslmode=disable"
	dbConfig := envGetter.GetDBConfig(true)
	dbString := dbConfig.DBString()
	if expected != dbString {
		t.Errorf("got %s\nwant %s \n", dbString, expected)
	}
}
