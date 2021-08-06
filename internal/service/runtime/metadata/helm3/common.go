package helm3

import (
	"encoding/json"
	"fmt"
	"strings"

	"baas/common/tools"

	"sigs.k8s.io/yaml"
)

// const 参数
const (
	RandStrLen = 6 // release name 随机字符串长度
)

// JSONValues values
type JSONValues interface {
	JSONMarshal() ([]byte, error)
}

// PeerReleaseName 创建 peer 节点的 release 名称
func PeerReleaseName(ns string) string {
	return fmt.Sprintf("peer-%s-%s", ns, tools.RandStringBytesMaskImprSrc(RandStrLen))
}

// PeerNamespace peer 节点的 namespace
func PeerNamespace(network, org string) string {
	return strings.ToLower(fmt.Sprintf("%s-%s", org, network))
}

// OrdererNamespace orderer 节点 namespace
func OrdererNamespace(network string) string {
	return strings.ToLower(fmt.Sprintf("%s", network))
}

// CreateNSRelease 创建 namespace 的 release 名称
func CreateNSRelease(ns string) string {
	return fmt.Sprintf("ns-%s", ns)
}

// CreateOrdererRelease 创建 orderer 节点的 release 名称
func CreateOrdererRelease(ns string) string {
	return fmt.Sprintf("orderer-%s-%s", ns, tools.RandStringBytesMaskImprSrc(RandStrLen))
}

//CreateEndpointRelease 创建 endpoint 的 release 名称
func CreateEndpointRelease(namespace string) string {
	return fmt.Sprintf("endpoints-%s-%s", namespace, tools.RandStringBytesMaskImprSrc(RandStrLen))
}

// PeerValues ...
type PeerValues struct {
	ConsortiumName  string   `json:"consortiumName"`
	PeerCount       string   `json:"peerCount"`
	OrgMspID        string   `json:"orgMspId"`
	OrgName         string   `json:"orgName"` // 组成 namespace （{orgName}-{ConsortiumName}）
	DbUsername      string   `json:"dbUsername"`
	DbPassword      string   `json:"dbPassword"`
	LogLevel        string   `json:"logLevel"`
	PreParams       []string `json:"preParams"`
	StartIndex      int      `json:"startIndex"`
	GossipBootStrap string   `json:"gossipBootStrap"`
	GossipLeader    bool     `json:"gossipLeader"`
	GossipElection  bool     `json:"gossipElection"`
	PeerTag         string   `json:"peerTag,omitempty"`
	ImageRepository string   `json:"imageRepository,omitempty"`
}

// TLSCollection ...
type TLSCollection struct {
	Cert []string `json:"cert,omitempty"`
	Key  string   `json:"key,omitempty"`
	CA   []string `json:"ca,omitempty"`
}

// MSPCollection ...
type MSPCollection struct {
	Admin []string `json:"admin,omitempty"`
	Sign  []string `json:"sign,omitempty"`
	CA    []string `json:"ca,omitempty"`
	Key   string   `json:"key,omitempty"`
}

// PodResourceChart ...
type PodResourceChart struct {
	CPU    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
}

// PreparedInfo ...
type PreparedInfo struct {
	NS              string            `json:"ns,omitempty"`
	Name            string            `json:"name,omitempty"`
	MSPID           string            `json:"mspID,omitempty"`
	LogLevel        string            `json:"logLevel,omitempty"`
	Request         *PodResourceChart `json:"request,omitempty"`
	Limit           *PodResourceChart `json:"limit,omitempty"`
	TLS             *TLSCollection    `json:"tls,omitempty"`
	MSP             *MSPCollection    `json:"msp,omitempty"`
	Genesis         string            `json:"genesis,omitempty"`
	GossipBootStrap string            `json:"gossipBootStrap,omitempty"`
	OUConfig        []string          `json:"ouconfig,omitempty"`
}

// ChartExt ...
type ChartExt struct {
	ImageRepository string `json:"imageRepository,omitempty"`
	ImageTag        string `json:"imageTag,omitempty"`
	ImagePullPolicy string `json:"imagePullPolicy,omitempty"`
	Network         string `json:"network,omitempty"`
	// CNCCGM          *NetSign       `json:"cnccgm,omitempty"`
	PreparedInfos   []PreparedInfo `json:"preparedInfos,omitempty"`
	GossipBootStrap string         `json:"gossipBootStrap,omitempty"`
	GossipLeader    bool           `json:"gossipLeader"`
	GossipElection  bool           `json:"gossipElection"`
	StorageName     string         `json:"storageName"`
	StorageSize     string         `json:"storageSize"`
	KubeConfig      string         `json:"kubeConfig"`
}

// JSONMarshal 序列化为 json 数据
func (ce *ChartExt) JSONMarshal() ([]byte, error) {
	valueDataByte, err := json.Marshal(ce)
	if err != nil {
		return nil, err
	}
	return valueDataByte, nil
}

func convertToMap(jj JSONValues) (map[string]interface{}, error) {
	var values = map[string]interface{}{}
	jBytes, err := jj.JSONMarshal()
	if err != nil {
		return nil, err
	}
	yBytes, err := yaml.JSONToYAML(jBytes)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yBytes, &values)
	if err != nil {
		return nil, err
	}
	return values, nil
}
