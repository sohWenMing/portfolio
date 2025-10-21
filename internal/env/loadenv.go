package loadenv

import (
	"encoding/base64"
	"os"
	"strings"

	"github.com/joho/godotenv"
	csrfProtect "github.com/sohWenMing/portfolio/internal/security/csrf_protect"
)

type EnvGetter struct{}

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
