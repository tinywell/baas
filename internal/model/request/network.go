package request

// NetCreate 网络创建请求
type NetCreate struct {
	NetInfo
}

// NetInit 网络初始化请求参数
type NetInit struct {
	Network       NetInfo       `json:"network,omitempty" `
	Runtime       string        `json:"runtime,omitempty" ` // Docker 、 k8s
	Members       []string      `json:"members,omitempty" `
	NodeOrderers  []OrgOrderers `json:"node_orderers,omitempty" `
	NodePeers     []OrgPeers    `json:"node_peers,omitempty" `
	NetSigns      []NetSign     `json:"net_signs,omitempty"`
	GenesisConfig OrdererConfig `json:"genesis_config,omitempty"`
	Version       string        `json:"version,omitempty" ` // fabric 版本（镜像版本）
}

// NetJoin 网络加入请求
type NetJoin struct {
}

// =======================================================================

// VMHost docker 宿主机
type VMHost struct {
	Hostname string      `json:"hostname,omitempty"`
	IP       string      `json:"ip,omitempty"`
	Port     int32       `json:"port,omitempty"`
	Type     int         `json:"type,omitempty"`
	Desc     string      `json:"desc,omitempty"`
	Config   interface{} `json:"config,omitempty"`
}

// OrdererConfig 共识配置
type OrdererConfig struct {
	Type              string `json:"type,omitempty"`
	BatchTimeout      int    `json:"batchTimeout,omitempty"` // s
	MaxMessageCount   int    `json:"maxMessageCount,omitempty"`
	AbsoluteMaxBytes  int    `json:"absoluteMaxBytes,omitempty"`  // M
	PreferredMaxBytes int    `json:"preferredMaxBytes,omitempty"` // KB
}

// NetSign 签名服务器信息
type NetSign struct {
	Name       string `json:"name,omitempty" `
	IP         string `json:"addr,omitempty" `
	Port       int    `json:"port,omitempty" `
	Password   string `json:"password,omitempty" `
	DataCenter string `json:"data_center,omitempty" `
}

// OrgOrderers 组织 orderer 节点信息
type OrgOrderers struct {
	MSPID string        `json:"mspid,omitempty" `
	Nodes []NodeOrderer `json:"nodes,omitempty" `
}

// OrgPeers 组织 peer 节点信息
type OrgPeers struct {
	MSPID string     `json:"mspid,omitempty" `
	Nodes []NodePeer `json:"nodes,omitempty" `
}

// NodePeer peer 节点信息
type NodePeer struct {
	NodeTempl
	IsBootstrap bool    `json:"is_bootstrap,omitempty" `
	IsAnchor    bool    `json:"is_anchor,omitempty" `
	StateDB     CouchDB `json:"state_db,omitempty" `
}

// NodeOrderer orderer 节点信息
type NodeOrderer struct {
	NodeTempl
}

// NodeTempl 节点通用信息
type NodeTempl struct {
	Name string `json:"name,omitempty" `
	NodeResource
}

// NodeResource 节点资源信息
type NodeResource struct {
	CPU        float64 `json:"cpu,omitempty"`     // 核
	Memory     float64 `json:"memory,omitempty" ` // G
	Stroge     int     `json:"stroge,omitempty"`  // G
	DataCenter string  `json:"data_center,omitempty"`
	HostID     int64   `json:"host_id,omitempty" `
}

// CouchDB CouchDB 状态数据库信息
type CouchDB struct {
	IP       string
	Port     int
	Admin    string
	Password string
}

// NetInfo 网络基本信息
type NetInfo struct {
	Name string `json:"name,omitempty"`
	Desc string `json:"desc,omitempty"`
	Type string `json:"type,omitempty"`
}

// DefOrderer 默认共识参数
var DefOrderer = OrdererConfig{
	Type:              "",
	BatchTimeout:      2,
	MaxMessageCount:   2000,
	AbsoluteMaxBytes:  99,
	PreferredMaxBytes: 20,
}
