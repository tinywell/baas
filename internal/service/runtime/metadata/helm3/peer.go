package helm3

import (
	"encoding/base64"
	"strings"

	"baas/internal/model"
	"baas/internal/service/runtime/metadata/common"
	"baas/pkg/runtime"
	"baas/pkg/runtime/helm3"
)

// charts
const (
	ChartPeerLevelDB = "baas-peer-plus-leveldb"
	// ChartPeerLevelDB = "baas-peer-plus-leveldb-new"
	ChartPeerCouchDB = "baas-peer-plus-couchdb"
)

// CreateData  peer 节点创建数据
func (dm *DataMachinePeer) CreateData(data *common.PeerData) runtime.ServiceMetadata {
	if data.Service.Runtime != model.RuntimeTypeHelm3 {
		return nil
	}

	svcData := helm3.NewInstallData()
	svcData.Name = data.Service.Name
	svcData.Namespace = PeerNamespace(data.NetworkName, data.Service.MSPID)
	svcData.ReleaseName = PeerReleaseName(svcData.Namespace)
	svcData.Chart = dm.preparePeerChart(data)
	svcData.Values = dm.preparePeerValues(data)
	return svcData
}

// DeleteData  peer 节点删除数据
func (dm *DataMachinePeer) DeleteData(data *common.PeerData) runtime.ServiceMetadata {
	return nil
}

func (dm *DataMachinePeer) preparePeerValues(data *common.PeerData) map[string]interface{} {
	//TODO: 部分信息需要结合配置参数设置
	ext := &ChartExt{
		ImageRepository: "",
		ImageTag:        data.Extra.Tag,
		ImagePullPolicy: "",
		Network:         data.NetworkName,
		PreparedInfos:   []PreparedInfo{dm.preparePeerInfo(data)},
		GossipLeader:    false,
		GossipElection:  true,
	}
	if len(data.BootStraps) > 0 {
		ext.GossipBootStrap = data.BootStraps
	}
	ce, err := convertToMap(ext)
	if err != nil {
		return nil
	}
	return ce
}

func (dm *DataMachinePeer) preparePeerInfo(data *common.PeerData) PreparedInfo {
	info := PreparedInfo{
		NS:       PeerNamespace(data.NetworkName, data.Service.MSPID),
		Name:     data.Service.Name,
		MSPID:    data.Service.MSPID,
		LogLevel: data.LogLevel,
		TLS: &TLSCollection{
			Cert: strings.Split(data.Extra.TLSCert, "\n"),
			Key:  base64.StdEncoding.EncodeToString([]byte(data.Extra.TLSKey)),
			CA:   strings.Split(data.Org.TLSCACert, "\n"),
		},
		MSP: &MSPCollection{
			Admin: strings.Split(data.Org.AdminCert, "\n"),
			Sign:  strings.Split(data.Extra.MSPCert, "\n"),
			CA:    strings.Split(data.Org.CACert, "\n"),
			Key:   base64.StdEncoding.EncodeToString([]byte(data.Extra.MSPKey)),
		},
		OUConfig: strings.Split(data.Extra.OUConfig, "\n"),
	}
	if len(data.BootStraps) > 0 {
		info.GossipBootStrap = data.BootStraps
	}
	return info
}

func (dm *DataMachinePeer) preparePeerChart(data *common.PeerData) string {
	chart := ""
	switch data.Extra.StateDB {
	case model.StateDBLevelDB:
		chart = ChartPeerLevelDB
	case model.StateDBCouchDB:
		chart = ChartPeerCouchDB
	default:
		chart = ChartPeerLevelDB
	}
	return chart
}
