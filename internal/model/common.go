package model

// BaaSData ...
type BaaSData struct {
	ID        int64 `json:"id,omitempty" db:"id"`
	TenantID  int64 `json:"tenant_id,omitempty" db:"tenant_id"`
	NetworkID int64 `json:"network_id,omitempty" db:"network_id"`
}

// HFNode hyperledger fabric node
type HFNode struct {
	MSPID    string
	MSPKey   string
	MSPCert  string
	TLSKey   string
	TLSCert  string
	OUConfig string
}

// 资源运行时类型
const (
	RuntimeTypeDocker = iota
	RuntimeTypeHelm2
	RuntimeTypeHelm3
	RuntimeTypeKubenetes
)

// 共识类型
const (
	OrdererTypeRaft = "etcdraft"
)
