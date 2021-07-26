package model

import (
	"encoding/json"
	"time"
)

// VMService 服务资源
type VMService struct {
	BaaSData
	MSPID      string                 `json:"mspid,omitempty" db:"mspid"`
	Name       string                 `json:"name,omitempty" db:"name"`
	Runtime    int                    `json:"runtime,omitempty" db:"runtime"`     // runtime 类型 : 0 - docker; 1 - k8s;
	LinkType   int                    `json:"link_type,omitempty" db:"link_type"` // 具体服务类型
	LinkID     int64                  `json:"link_id,omitempty" db:"link_id"`     // 关联具体资源 ID
	CFGRaw     []byte                 `json:"cfg_raw,omitempty" db:"cfg"`         // 具体配置信息
	CFG        map[string]interface{} `json:"cfg,omitempty" db:"-"`
	Status     int                    `json:"status,omitempty" db:"status"`
	DataCenter string                 `json:"data_center,omitempty" db:"data_center"`
	DCID       int64
	DCMetadata []byte    `json:"dc_metadata,omitempty" db:"-"`
	CreateTime time.Time `json:"create_time,omitempty" db:"create_time"`
	UpdateTime time.Time `json:"update_time,omitempty" db:"update_time"`
	Creator    string    `json:"creator,omitempty" db:"creator"`
}

// DataCenterDocker ...
type DataCenterDocker struct {
	Name    string `json:"name,omitempty" db:"name"`
	Host    string `json:"host,omitempty" db:"host"`
	Port    int    `json:"port,omitempty" db:"port"`
	TLS     bool   `json:"tls,omitempty" db:"tls"`
	TLSCert string `json:"tls_cert,omitempty" db:"tls_cert"` // tls cert path
	TLSKey  string `json:"tls_key,omitempty" db:"tls_key"`   // tls key path
	TLSCA   string `json:"tlsca,omitempty" db:"tlsca"`       // tls root cert path
	Sock    string `json:"sock,omitempty" db:"sock"`         // .sock path
}

// ToBytes ...
func (dc *DataCenterDocker) ToBytes() ([]byte, error) {
	return json.Marshal(dc)
}

// FromBytes ...
func (dc *DataCenterDocker) FromBytes(raw []byte) error {
	return json.Unmarshal(raw, dc)
}
