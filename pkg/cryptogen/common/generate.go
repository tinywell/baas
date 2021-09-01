package common

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"net"
	"time"

	"github.com/hyperledger/fabric/bccsp"
	"github.com/pkg/errors"
)

// KeyGenerater 密钥生成器
type KeyGenerater interface {
	GenKey(keystore string) (bccsp.Key, error)
}

// GenerateKeyPair 生成密钥对
func GenerateKeyPair(keystore string, gen KeyGenerater) (bccsp.Key, bccsp.Key, error) {
	prikey, err := gen.GenKey(keystore)
	if err != nil {
		return nil, nil, errors.WithMessage(err, "生成私钥失败")
	}
	pubkey, err := prikey.PublicKey()
	if err != nil {
		return nil, nil, errors.WithMessage(err, "获取公钥失败")
	}

	return prikey, pubkey, nil
}

// GenerateCert 生成证书
func GenerateCert(prikey bccsp.Key, pubkey interface{}, signCA crypto.Signer, parent *Cert, spec *NodeSpec, t CertType) (cert string, err error) {
	template := certificate(spec)
	template.SubjectKeyId = prikey.SKI()

	var tmpCert *Cert
	if parent == nil { // ca
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageDigitalSignature |
			x509.KeyUsageKeyEncipherment | x509.KeyUsageCertSign |
			x509.KeyUsageCRLSign
		template.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageAny}
	} else {
		template.KeyUsage = x509.KeyUsageDigitalSignature

	}
	tmpCert, _ = NewCertFromX509Temp(&template)
	ncert, err := CreateCertificate(rand.Reader, tmpCert, parent, pubkey, signCA, t)
	if err != nil {
		return "", err
	}

	certPem, err := ncert.ToPem()
	if err != nil {
		return cert, errors.WithMessage(err, "证书转 pem 格式失败")
	}
	return string(certPem), nil
}

func certificate(spec *NodeSpec) x509.Certificate {

	//basic template to use
	template := x509Template()

	//set the organization for the subject
	subject := subjectTemplateAdditional(spec.Country,
		spec.Province,
		spec.Locality,
		spec.OrganizationalUnit,
		spec.StreetAddress,
		spec.PostalCode)
	if len(spec.Organization) > 0 {
		subject.Organization = []string{spec.Organization}
	}
	subject.CommonName = spec.CommonName

	template.Subject = subject
	for _, san := range spec.SANS {
		// try to parse as an IP address first
		ip := net.ParseIP(san)
		if ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, san)
		}
	}
	return template
}

// default template for X509 certificates
func x509Template() x509.Certificate {

	// generate a serial number
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, serialNumberLimit)

	// set expiry to around 10 years
	expiry := 3650 * 24 * time.Hour
	// backdate 5 min
	notBefore := time.Now().Add(-5 * time.Minute).UTC()

	//basic template to use
	x509 := x509.Certificate{
		SerialNumber:          serialNumber,
		NotBefore:             notBefore,
		NotAfter:              notBefore.Add(expiry).UTC(),
		BasicConstraintsValid: true,
	}
	return x509

}

// Additional for X509 subject
func subjectTemplateAdditional(country, province, locality, orgUnit, streetAddress, postalCode string) pkix.Name {
	name := subjectTemplate()
	if len(country) >= 1 {
		name.Country = []string{country}
	}
	if len(province) >= 1 {
		name.Province = []string{province}
	}

	if len(locality) >= 1 {
		name.Locality = []string{locality}
	}
	if len(orgUnit) >= 1 {
		name.OrganizationalUnit = []string{orgUnit}
	}
	if len(streetAddress) >= 1 {
		name.StreetAddress = []string{streetAddress}
	}
	if len(postalCode) >= 1 {
		name.PostalCode = []string{postalCode}
	}
	return name
}

// default template for X509 subject
func subjectTemplate() pkix.Name {
	return pkix.Name{
		Country:  []string{"CN"},
		Locality: []string{"Beijing"},
		Province: []string{"Beijing"},
	}
}
