package model

// BaaSData ...
type BaaSData struct {
	ID        int64 `json:"id,omitempty" db:"id"`
	TenantID  int64 `json:"tenant_id,omitempty" db:"tenant_id"`
	NetworkID int64 `json:"network_id,omitempty" db:"network_id"`
}

// HFNode hyperledger fabric node
type HFNode struct {
	MSPID    string `json:"mspid,omitempty" db:"mspid"`
	Name     string `json:"name,omitempty" db:"name"`
	MSPKey   string `json:"msp_key,omitempty" db:"msp_key"`
	MSPCert  string `json:"msp_cert,omitempty" db:"msp_cert"`
	TLSKey   string `json:"tls_key,omitempty" db:"tls_key"`
	TLSCert  string `json:"tls_cert,omitempty" db:"tls_cert"`
	OUConfig string `json:"ou_config,omitempty" db:"ou_config"`
}

// 资源运行时类型
const (
	RuntimeTypeDocker = iota
	RuntimeTypeHelm2
	RuntimeTypeHelm3
	RuntimeTypeKubenetes
)

// RuntimeTypeNameDocker 运行时类型（名称）
const (
	RuntimeTypeNameDocker    = "docker"
	RuntimeTypeNameHelm2     = "helm2"
	RuntimeTypeNameHelm3     = "helm3"
	RuntimeTypeNameKubenetes = "kubenetes"
)

// runtime 转换
var (
	RuntimeTypeNameValue = map[string]int{
		RuntimeTypeNameDocker:    RuntimeTypeDocker,
		RuntimeTypeNameHelm2:     RuntimeTypeHelm2,
		RuntimeTypeNameHelm3:     RuntimeTypeHelm3,
		RuntimeTypeNameKubenetes: RuntimeTypeKubenetes,
	}

	RuntimeTypeValueName = map[int]string{
		RuntimeTypeDocker:    RuntimeTypeNameDocker,
		RuntimeTypeHelm2:     RuntimeTypeNameHelm2,
		RuntimeTypeHelm3:     RuntimeTypeNameHelm3,
		RuntimeTypeKubenetes: RuntimeTypeNameKubenetes,
	}
)

// VMService 类型
const (
	VMServiceTypePeer = iota
	VMServiceTypeOrderer
	VMServiceTypeHost
)

// 共识类型
const (
	OrdererTypeRaft = "etcdraft"
	OrdererTypeSolo = "solo"
)

// 加解密套件类型
const (
	CryptoTypeSW = "SW"
	CryptoTypeGM = "GM"
)
