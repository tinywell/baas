package helm3

import (
	"baas/internal/model"
	module "baas/internal/model"
	"encoding/base64"
	"strings"

	"baas/internal/service/runtime/metadata/common"
	"baas/pkg/runtime"
	"baas/pkg/runtime/helm3"
)

// charts
const (
	ChartOrdererRaft = "baas-orderer-raft"
	ChartOrdererSolo = "baas-orderer-solo"
)

// CreateData 生成创建 orderer 节点的 helm 模板数据
func (dm *DataMachineOrderer) CreateData(data *common.OrdererData) runtime.ServiceMetadata {
	if data.Service.Runtime != module.RuntimeTypeHelm3 {
		return nil
	}
	//TODO: namespace 如何处理？
	svcData := helm3.NewInstallData()
	svcData.Name = data.Service.Name
	svcData.Namespace = OrdererNamespace(data.NetworkName)
	svcData.ReleaseName = OrdererReleaseName(svcData.Namespace)
	svcData.Chart = dm.prepareChart(data)
	svcData.Values = dm.prepareValues(data)
	return svcData
}

// DeleteData 生成删除 ordrerer 节点的 helm 模板数据
func (dm *DataMachineOrderer) DeleteData(data *common.OrdererData) runtime.ServiceMetadata {
	panic("not implemented") // TODO: Implement
}

func (dm *DataMachineOrderer) prepareValues(data *common.OrdererData) map[string]interface{} {
	ext := &ChartExt{
		ImageRepository: "",
		ImageTag:        data.Extra.Tag,
		ImagePullPolicy: "",
		Network:         data.NetworkName,
		PreparedInfos:   []PreparedInfo{dm.prepareInfo(data)},
	}
	ce, err := convertToMap(ext)
	if err != nil {
		return nil
	}
	return ce
}

func (dm *DataMachineOrderer) prepareInfo(data *common.OrdererData) PreparedInfo {
	info := PreparedInfo{
		NS:       OrdererNamespace(data.NetworkName),
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
		Genesis:  base64.StdEncoding.EncodeToString(data.Genesis),
	}
	return info
}
func (dm *DataMachineOrderer) prepareChart(data *common.OrdererData) string {
	chart := ""
	switch data.OrdererType {
	case model.OrdererTypeRaft:
		chart = ChartOrdererRaft
	case model.OrdererTypeSolo:
		chart = ChartOrdererSolo
	}
	return chart
}
