package csfr_protect

import (
	"net/http"

	"github.com/gorilla/csrf"
)

type CSRFKeyData struct {
	KeyString string
	KeyBytes  []byte
}

type CsrfEnvGetter interface {
	GetCSRFKey() CSRFKeyData
}

func GetCSRFKey(c CsrfEnvGetter) []byte {
	return c.GetCSRFKey().KeyBytes
}

func LoadCSRFMW(envPath string, envGetter CsrfEnvGetter) func(next http.Handler) http.Handler {
	csrf := csrf.Protect(
		GetCSRFKey(envGetter),
		csrf.TrustedOrigins(
			[]string{
				"localhost:3000",
				"128.199.244.100",
			},
		),
		csrf.Secure(false),
	)
	return csrf
}
