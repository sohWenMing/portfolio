package databaseintegrationtesting

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	dbinfra "github.com/sohWenMing/portfolio/internal/db_infra"
	loadenv "github.com/sohWenMing/portfolio/internal/env"
)

var envGetter *loadenv.EnvGetter

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

func TestInitDB(t *testing.T) {
	db, err := dbinfra.InitDB(envGetter.GetDBConfig(true))
	if err != nil {
		t.Errorf("didn't expect error, got %v", err)
		return
	}
	err = db.Close()
	if err != nil {
		t.Errorf("didn't expect error, got %v", err)
		return
	}
}
