package common

import (
	"crypto"

	"github.com/hyperledger/fabric/bccsp"
)

// crypto type
const (
	CryptoTypeGM   = "GM"
	CryptoTypeSW   = "SW"
	CryptoTypeSWGM = "SWGM" // 软国密
)

// CN
const (
	Country  = "CN"
	Province = "BeiJing"
	Locality = "BeiJing"
)

// cert 类型
const (
	CAType    = "ca"
	TLSCAType = "tlsca"
)

// NodeOU
const (
	CLIENT = iota
	ORDERER
	PEER
	ADMIN
)

// NodeOU
const (
	CLIENTOU  = "client"
	PEEROU    = "peer"
	ADMINOU   = "admin"
	ORDEREROU = "orderer"
)

var nodeOUMap = map[int]string{
	CLIENT:  CLIENTOU,
	PEER:    PEEROU,
	ADMIN:   ADMINOU,
	ORDERER: ORDEREROU,
}

// ToolSet 工具集
type ToolSet struct {
	CryptoType    string
	GenKey        FcnKeyGen
	ToPem         FcnKeyToPem
	PubToPem      FcnCertToPem
	CreateSigner  FcnCreateSigner
	GetRootSigner FcnGetRootSigner
}

// FcnKeyGen 密钥生成
type FcnKeyGen func(keystore string) (bccsp.Key, error)

// FcnKeyToPem 密钥 pem 格式
type FcnKeyToPem func(keystore string) ([]byte, error)

// FcnCertToPem 公钥证书 pem 格式
type FcnCertToPem func(raw []byte) ([]byte, error)

// FcnCreateSigner 通过私钥生成签名 signer
type FcnCreateSigner func(prikey bccsp.Key) (crypto.Signer, error)

// FcnGetRootSigner 根证书签名用 signer
type FcnGetRootSigner func() (crypto.Signer, error)

// Member 组织成员证书信息（用户或节点）
type Member struct {
	Name    string
	TLSCert string
	TLSKey  string
	MSPCert string
	MSPKey  string
}

// Organization 组织证书信息
type Organization struct {
	Name      string // req -> domain orderer.citic.com ; peer.citit.com
	MSPID     string
	TLSCACert string
	TLSCAKey  string
	MSPCACert string
	MSPCAKey  string
	Unit      string
}

// OrgSpec 组织基本信息
type OrgSpec struct {
	Name          string   `json:"Name"`
	CommonName    string   // 组织域名
	EnableNodeOUs bool     `json:"EnableNodeOUs"` // 是否生成 节点下的msp的config.yaml配置文件
	CA            NodeSpec `json:"CA"`
	//Template      NodeTemplate `json:"Template"`
	//Specs         []NodeSpec   `json:"Specs"`
	//Users UsersSpec
}

// NewOrgSpec 返回新的组织基本信息
func NewOrgSpec() *OrgSpec {
	orgspec := &OrgSpec{}
	return orgspec
}

// NodeSpec 节点基本信息
type NodeSpec struct {
	Organization       string // 组织域名
	CommonName         string // 节点域名
	Country            string // 国家
	Province           string // 省份
	Locality           string // 地区
	OrganizationalUnit string // OU ('admin','peer','orderer','client')
	StreetAddress      string
	PostalCode         string
	SANS               []string
}

// UsersSpec 用户基本信息
type UsersSpec struct {
	Count int `json:"Count"`
}
