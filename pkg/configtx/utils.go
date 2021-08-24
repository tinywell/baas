package configtx

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric-config/configtx"
)

func transCert(cert string) (*x509.Certificate, error) {
	p, _ := pem.Decode([]byte(cert))
	if p == nil {
		return nil, errors.New("证书 pem 格式错误")
	}
	return x509.ParseCertificate(p.Bytes)
}

func getOrgPolicy(orgMspID string, mems []string) configtx.Policy {
	member := make([]string, 0, len(mems))
	for _, m := range mems {
		member = append(member, fmt.Sprintf("'%s.%s'", orgMspID, m))
	}
	return configtx.Policy{
		Type: configtx.SignaturePolicyType,
		Rule: fmt.Sprintf("OR(%s)", strings.Join(member, ",")),
	}
}
