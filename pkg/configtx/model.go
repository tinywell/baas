package configtx

import "time"

// 配置中的常量
const (
	ConsortiumName = "SampleConsortium"
	HashFamily     = "SHA2"
	HashFunc       = "SHA256"
	ModPolicy      = "Admin"
)

// Organization 组织信息
type Organization struct {
	MSPID     string
	MSPCert   string
	TLSCert   string
	AdminCert string
	Endpoints []Endpoint
	NodeOU    bool
}

// Endpoint .
type Endpoint string

// OrdererConfig 共识相关配置
type OrdererConfig struct {
	OdererType string
	Cutter     CutterConfig
	Raft       RaftConsentor
}

// RaftConsentor raft 节点配置
type RaftConsentor struct {
	Address    string
	ServerCert string
	ClientCert string
}

// CutterConfig 出块配置
type CutterConfig struct {
	BatchTimeout      time.Duration
	BatchSizeAbsolute uint32
	BatchSizePrefer   uint32
	BatchSizeMaxCount uint32
}
