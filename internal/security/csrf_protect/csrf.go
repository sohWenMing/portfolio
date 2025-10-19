package csfr_protect

type CSRFKeyData struct {
	KeyString string
	KeyBytes  []byte
}

type csrfEnvGetter interface {
	GetCSRFKey() CSRFKeyData
}

func GetCSRFKey(c csrfEnvGetter) []byte {
	return c.GetCSRFKey().KeyBytes
}
