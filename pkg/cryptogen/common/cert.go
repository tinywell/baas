package common

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
)

// CertType 证书类型
type CertType int

// 证书类型
const (
	CertTypeX509 CertType = iota
	CertTypeGM
)

// Cert 证书抽象类型
type Cert struct {
	certType  CertType
	certInner []byte
}

// NewCertFromBytes 基于证书类型及其序列化数据生成 Cert 实例
func NewCertFromBytes(t CertType, raw []byte) *Cert {
	return &Cert{
		certType:  t,
		certInner: raw,
	}
}

// NewCertFromPem 基于证书类型及其 Pem 格式序列化数据生成 Cert 实例
func NewCertFromPem(t CertType, raw []byte) (*Cert, error) {
	block, _ := pem.Decode(raw)
	if block == nil {
		return nil, errors.New("pem 证书格式错误")
	}
	return &Cert{
		certType:  t,
		certInner: block.Bytes,
	}, nil
}

// NewCertFromX509Cert 基于 x509.Certificate 实例生成 Cert 实例
func NewCertFromX509Cert(cert *x509.Certificate) (*Cert, error) {
	return &Cert{
		certType:  CertTypeX509,
		certInner: cert.Raw,
	}, nil
}

// Type 获取 Cert 原证书类型
func (c *Cert) Type() CertType {
	return c.certType
}

// ToX509 转化为 x509.Certificate 实例
func (c *Cert) ToX509() (*x509.Certificate, error) {
	return x509.ParseCertificate(c.certInner)
}

// ToPem 转化为 pem 格式序列化数据
func (c *Cert) ToPem() ([]byte, error) {
	block := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: c.certInner,
	}
	return pem.EncodeToMemory(block), nil
}

// CreateCertificate 签发证书
func CreateCertificate(rand io.Reader, template, parent *Cert, pub, priv interface{}, t CertType) (cert *Cert, err error) {
	cert.certType = t
	switch t {
	case CertTypeX509:
		p, err := parent.ToX509()
		if err != nil {
			return nil, err
		}
		temp, err := template.ToX509()
		if err != nil {
			return nil, err
		}
		c, err := x509.CreateCertificate(rand, temp, p, pub, priv)
		if err != nil {
			return nil, err
		}
		cert.certInner = c
	case CertTypeGM:
		//TODO:
	}
	return
}
