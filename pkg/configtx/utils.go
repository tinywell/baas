package configtx

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func transCert(cert string) (*x509.Certificate, error) {
	p, _ := pem.Decode([]byte(cert))
	if p == nil {
		return nil, errors.New("证书 pem 格式错误")
	}
	return x509.ParseCertificate(p.Bytes)
}
