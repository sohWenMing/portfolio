package loadenv

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	csrfProtect "github.com/sohWenMing/portfolio/internal/security/csrf_protect"
)

type DBConfig struct {
	DBPort         string
	DBScheme       string
	DBHost         string
	DBUser         string
	DBPassword     string
	DBDatabasename string
	DBSSLmode      string
}

func (dbConfig *DBConfig) DBString() string {
	return fmt.Sprintf(
		"host=%s dbname=%s user=%s password=%s sslmode=%s",
		dbConfig.DBHost,
		dbConfig.DBDatabasename,
		dbConfig.DBUser,
		dbConfig.DBPassword,
		dbConfig.DBSSLmode,
	)
}

// host=localhost dbname=mydatabase user=myuser password=mypassword sslmode=disable
type EnvGetter struct{}

func (*EnvGetter) GetDBConfig(isTestLocalConn bool) *DBConfig {
	dbConfig :=
		&DBConfig{
			os.Getenv("DBPORT"),
			os.Getenv("DBSCHEME"),
			os.Getenv("DBHOST"),
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
			os.Getenv("POSTGRES_DB"),
			os.Getenv("DBSSLMODE"),
		}
	if isTestLocalConn {
		dbConfig.DBHost = "localhost"
	}
	return dbConfig
}

// gets the secret key that is used for the generation of the CSRF token
func (*EnvGetter) GetCSRFKey() csrfProtect.CSRFKeyData {
	csrfKey := os.Getenv("CSRFKEY")
	decoded, err := base64.StdEncoding.DecodeString(csrfKey)
	if err != nil {
		panic(err)
	}
	return csrfProtect.CSRFKeyData{
		KeyString: csrfKey,
		KeyBytes:  decoded,
	}
}

// gets the trust origins that are used for the CSRF token
func (*EnvGetter) GetTrustedOrigins() []string {
	trustedOrigins := os.Getenv("CSRFTRUSTEDORIGINS")
	trustedOriginsSlice := strings.Split(trustedOrigins, "|")
	return trustedOriginsSlice
}

func LoadEnv(envPath string) (getter *EnvGetter, err error) {
	err = godotenv.Load(envPath)
	if err != nil {
		return nil, err
	}
	return &EnvGetter{}, nil
}
