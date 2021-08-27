package configtx

import (
	"time"

	"github.com/hyperledger/fabric-config/configtx"
)

// 配置中的常量
const (
	ConsortiumName = "SampleConsortium"
	HashFamily     = "SHA2"
	HashFunc       = "SHA256"
	ModPolicyAdmin = "Admin"
)

// 默认 orderer 配置
const (
	DefBatchTimeout              = 2 * time.Second
	DefOrdererType               = "etcdraft"
	DefBatchSizeMaxCount  uint32 = 1000
	DefBatchSizeAbsolute  uint32 = 99 * 1024 * 1024 // 99M
	DefBatchSizePerffered uint32 = 20 * 1024 * 1024 // 20M
)

// raft 默认配置
const (
	DefTickInterval         string = ""
	DefElectionTick         uint32 = 0
	DefHeartbeatTick        uint32 = 0
	DefMaxInflightBlocks    uint32 = 0
	DefSnapshotIntervalSize uint32 = 0
)

// NodeOU
const (
	NodeOUPeer   = "peer"
	NodeOUClient = "client"
	NodeOUAdmin  = "admin"
	NodeOUMember = "member"
)

// 默认策略
var (
	DefReaderPolicy = configtx.Policy{
		Type: configtx.ImplicitMetaPolicyType,
		Rule: "ANY Readers",
	}
	DefWriterPolicy = configtx.Policy{
		Type: configtx.ImplicitMetaPolicyType,
		Rule: "ANY Writers",
	}
	DefAdminPolicy = configtx.Policy{
		Type: configtx.ImplicitMetaPolicyType,
		Rule: "MAJORITY Admins",
	}
	DefEndorsePolicy = configtx.Policy{
		Type: configtx.ImplicitMetaPolicyType,
		Rule: "MAJORITY Endorsement",
	}
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
	Raft       []*RaftConsentor
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
