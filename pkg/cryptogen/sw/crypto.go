package sw

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/hyperledger/fabric/bccsp"
	"github.com/hyperledger/fabric/bccsp/factory"
	"github.com/hyperledger/fabric/bccsp/signer"
	"github.com/pkg/errors"

	"baas/pkg/cryptogen/common"
)

var (
	csp bccsp.BCCSP
)

// Gen 原生 sw 证书颁发
type Gen struct {
}

// GenerateOrgCA 颁发组织根证书
func (g *Gen) GenerateOrgCA(spec *common.NodeSpec) (org common.Organization, err error) {
	org.Name = spec.CommonName
	capath, tlscapath, err := common.GenerateCaDir(spec.Organization)
	if err != nil {
		return
	}
	defer func() {
		common.Cleanup(capath)
		common.Cleanup(tlscapath)
	}()
	spec.CommonName = fmt.Sprintf("ca.%s", spec.Organization)
	key, cert, err := g.GenerateCACertPair(capath, spec)
	if err != nil {
		return
	}
	org.MSPCACert = cert
	org.MSPCAKey = key

	spec.CommonName = fmt.Sprintf("tlsca.%s", spec.Organization)
	tlskey, tlscert, err := g.GenerateCACertPair(tlscapath, spec)
	if err != nil {
		return
	}
	org.TLSCACert = tlscert
	org.TLSCAKey = tlskey
	return
}

// GenarateMember 颁发组织成员证书
func (g *Gen) GenarateMember(spec *common.NodeSpec, CA *common.Organization) (member common.Member, err error) {
	member.Name = spec.CommonName
	mempath, tlsmempath, err := common.GenerateMemberDir(spec.Organization, spec.CommonName)
	if err != nil {
		return
	}
	defer func() {
		common.Cleanup(mempath)
		common.Cleanup(tlsmempath)
	}()
	spec.Organization = ""
	key, cert, err := g.GenerateMemberCertPair(mempath, spec, CA.MSPCAKey, CA.MSPCACert)
	if err != nil {
		return
	}
	member.MSPKey = key
	member.MSPCert = cert

	spec.OrganizationalUnit = ""
	spec.SANS = append(spec.SANS, spec.CommonName)
	tlskey, tlscert, err := g.GenerateMemberCertPair(tlsmempath, spec, CA.TLSCAKey, CA.TLSCACert)
	if err != nil {
		return
	}
	member.TLSKey = tlskey
	member.TLSCert = tlscert
	return
}

// GenKey impl for common.KeyGenerater
func (g *Gen) GenKey(keystore string) (bccsp.Key, error) {
	opts := &factory.FactoryOpts{
		ProviderName: "SW",
		SwOpts: &factory.SwOpts{
			HashFamily: "SHA2",
			SecLevel:   256,
			FileKeystore: &factory.FileKeystoreOpts{
				KeyStorePath: keystore,
			},
		},
	}

	swcsp, err := factory.GetBCCSPFromOpts(opts)
	if err != nil {
		return nil, err
	}
	csp = swcsp
	return csp.KeyGen(&bccsp.ECDSAP256KeyGenOpts{Temporary: false})
}

// GenerateCACertPair ..
func (g *Gen) GenerateCACertPair(keystore string, spec *common.NodeSpec) (key, cert string, err error) {
	return g.GenerateCertPair(keystore, spec, nil, nil)
}

// GenerateMemberCertPair ..
func (g *Gen) GenerateMemberCertPair(keystore string, spec *common.NodeSpec, caKey, caCert string) (key, cert string, err error) {

	capri, err := LoadECPrivateKey([]byte(caKey))
	if err != nil {
		return
	}

	signCA, err := signer.New(csp, capri)
	if err != nil {
		return
	}

	parent, err := common.NewCertFromPem(common.CertTypeX509, []byte(caCert))
	if err != nil {
		return
	}
	return g.GenerateCertPair(keystore, spec, signCA, parent)
}

// GenerateCertPair ..
func (g *Gen) GenerateCertPair(keystore string, spec *common.NodeSpec, signCA crypto.Signer, parent *common.Cert) (key, cert string, err error) {

	pri, pub, err := common.GenerateKeyPair(keystore, g)
	if err != nil {
		err = errors.WithMessage(err, "生成密钥对失败")
		return
	}

	pubkey, err := GetECPublicKey(pub)
	if err != nil {
		err = errors.WithMessage(err, "获取 EC 公钥失败")
		return
	}

	if signCA == nil {
		signCA, err = signer.New(csp, pri)
		if err != nil {
			err = errors.WithMessage(err, "从私钥生成 signer 失败")
			return
		}
	}

	cert, err = common.GenerateCert(pri, pubkey, signCA, parent, spec, common.CertTypeX509)
	if err != nil {
		err = errors.WithMessage(err, "签发证书失败")
		return
	}
	pripem, err := common.LoadPrivateKey(keystore)
	if err != nil {
		err = errors.WithMessage(err, "提取私钥数据失败")
		return
	}
	key = string(pripem)
	return
}

// GetECPublicKey .
func GetECPublicKey(key bccsp.Key) (*ecdsa.PublicKey, error) {
	// get the public key
	pubKey, err := key.PublicKey()
	if err != nil {
		return nil, err
	}
	pubKeyBytes, err := pubKey.Bytes()
	if err != nil {
		return nil, err
	}
	// unmarshal using pkix
	ecPubKey, err := x509.ParsePKIXPublicKey(pubKeyBytes)
	if err != nil {
		return nil, err
	}
	return ecPubKey.(*ecdsa.PublicKey), nil
}

// LoadECPrivateKey ..
func LoadECPrivateKey(raw []byte) (bccsp.Key, error) {
	block, _ := pem.Decode(raw)
	if block == nil {
		return nil, errors.Errorf("pem 密钥格式错误")
	}

	priv, err := csp.KeyImport(block.Bytes, &bccsp.ECDSAPrivateKeyImportOpts{Temporary: true})
	if err != nil {
		return nil, errors.WithMessage(err, "导入密钥失败")
	}
	return priv, nil
}
