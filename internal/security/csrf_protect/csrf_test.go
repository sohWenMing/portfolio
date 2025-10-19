package csfr_protect

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	exitCode := m.Run()
	os.Exit(exitCode)
}
