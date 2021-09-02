package request

// NetCreate 网络创建请求
type NetCreate struct {
	NetInfo
}

// NetInit 网络初始化请求参数
type NetInit struct {
	Network       NetInfo                  `json:"network,omitempty" db:"network"`
	Runtime       string                   `json:"runtime,omitempty" db:"runtime"` // Docker 、 k8s
	Members       []Org                    `json:"members,omitempty" db:"members"`
	NodeOrderers  map[string][]NodeOrderer `json:"node_orderers,omitempty" db:"node_orderers"`
	NodePeers     map[string][]NodePeer    `json:"node_peers,omitempty" db:"node_peers"`
	ImagePeer     string                   `json:"image_peer,omitempty" db:"image_peer"`
	ImagesOrderer string                   `json:"images_orderer,omitempty" db:"images_orderer"`
	StateDB       string                   `json:"state_db,omitempty" db:"state_db"`       //levelDB、couchDB
	CryptoType    string                   `json:"crypto_type,omitempty" db:"crypto_type"` //SW、GM
	NetSigns      []NetSign                `json:"net_signs,omitempty" db:"net_signs"`
	Hosts         []RuntimeHost            `json:"hosts,omitempty" db:"hosts"`
	GenesisConfig OrdererConfig            `json:"genesis_config,omitempty" db:"genesis_config"`
	Version       string                   `json:"version,omitempty" db:"version"` // fabric 版本（镜像版本）
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

// RuntimeHost 运行时配置
type RuntimeHost struct {
	Name      string
	Host      string
	Scheme    string
	Type      string
	Desc      string
	TLS       bool
	TLSConfig struct {
		TLSKey  string
		TLSCert string
		TLSCA   string
	}
	HelmConfig struct {
		RepoConfig struct {
			Repo     string
			Private  bool
			Username string
			Password string
		}
		Kubefile string
	}
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

// Org 组织基本信息
type Org struct {
	MSPID  string `json:"mspid" `
	Domain string `json:"domain" `
}

// NodePeer peer 节点信息
type NodePeer struct {
	NodeTempl
	Port        int
	IsBootstrap bool    `json:"is_bootstrap,omitempty" `
	IsAnchor    bool    `json:"is_anchor,omitempty" `
	StateDB     CouchDB `json:"state_db,omitempty" `
}

// NodeOrderer orderer 节点信息
type NodeOrderer struct {
	NodeTempl
	Port int `json:"port,omitempty" db:"port"`
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
	Type:              "etcdraft",
	BatchTimeout:      2,
	MaxMessageCount:   2000,
	AbsoluteMaxBytes:  99,
	PreferredMaxBytes: 20,
}
