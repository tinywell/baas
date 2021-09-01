package docker

import (
	"baas/common/tools"
	"baas/internal/model"
	"baas/internal/service/runtime/metadata/common"
	"baas/pkg/runtime"
	"baas/pkg/runtime/docker"
	"encoding/base64"
	"strconv"
)

// orderer 节点配置常量
const (
	PATHGenesis  = "/etc/hyperledger/orderer/"
	GenesisFile  = "orderer.genesis.block"
	ImageOrderer = "hyperledger/fabric-orderer"
	PortOrderer  = 7050
)

// 默认值
var (
	DefaultEnvOrderer = map[string]string{
		"FABRIC_LOGGING_SPEC":                       "INFO",
		"ORDERER_GENERAL_LISTENADDRESS":             "0.0.0.0",
		"ORDERER_GENERAL_LISTENPORT":                "7050",
		"ORDERER_GENERAL_GENESISMETHOD":             "file",
		"ORDERER_GENERAL_GENESISFILE":               PATHGenesis + GenesisFile,
		"ORDERER_GENERAL_LOCALMSPID":                "OrdererMSP",
		"ORDERER_GENERAL_LOCALMSPDIR":               PATHMSP,
		"ORDERER_OPERATIONS_LISTENADDRESS":          "0.0.0.0:17050",
		"ORDERER_GENERAL_TLS_ENABLED":               "true",
		"ORDERER_GENERAL_TLS_PRIVATEKEY":            PATHTLSKey,
		"ORDERER_GENERAL_TLS_CERTIFICATE":           PATHTLSCert,
		"ORDERER_GENERAL_TLS_ROOTCAS":               "[" + PATHTLSCA + "]",
		"ORDERER_KAFKA_TOPIC_REPLICATIONFACTOR":     "1",
		"ORDERER_KAFKA_VERBOSE":                     "true",
		"ORDERER_GENERAL_CLUSTER_CLIENTCERTIFICATE": PATHTLSCert,
		"ORDERER_GENERAL_CLUSTER_CLIENTPRIVATEKEY":  PATHTLSKey,
		"ORDERER_GENERAL_CLUSTER_ROOTCAS":           "[" + PATHTLSCA + "]",
	}
)

// CreateData orderer 节点启动 docker 相关参数
func (dm *DataMachineOrderer) CreateData(data *common.OrdererData) runtime.ServiceMetadata {
	if data.Service.Runtime != model.RuntimeTypeDocker {
		return nil
	}
	svcData := docker.NewSingleServiceData()
	svcData.Name = data.Service.Name
	if len(data.Extra.Image) == 0 {
		data.Extra.Image = ImageOrderer
	}
	if len(data.Extra.Tag) > 0 {
		svcData.Image = data.Extra.Image + ":" + data.Extra.Tag
	} else {
		svcData.Image = data.Extra.Image
	}
	svcData.Ports = dm.preparePorts(data)
	// svcData.Volumes = prepareVolumes(data)
	svcData.ENVs = dm.prepareEnvs(data)
	svcData.ExtraHosts = data.ExtraHost
	svcData.Network = prepareNetwork(data.NetworkName, data.Extra.MSPID)
	svcData.CMDs = dm.prepareCMDs(data)

	return svcData
}

// DeleteData orderer 节点删除 docker 相关参数
func (dm *DataMachineOrderer) DeleteData(data *common.OrdererData) runtime.ServiceMetadata {
	panic("not implemented") // TODO: Implement
}

func (dm *DataMachineOrderer) prepareEnvs(data *common.OrdererData) []string {
	envs := tools.CopyStrMap(DefaultEnvOrderer)

	envs["ORDERER_GENERAL_LOCALMSPID"] = data.Extra.MSPID
	envs["ORDERER_GENERAL_LISTENPORT"] = strconv.Itoa(data.Extra.Port)

	envStr := make([]string, 0, len(envs))
	for k, v := range envs {
		envStr = append(envStr, k+"="+v)
	}
	return envStr
}

func (dm *DataMachineOrderer) preparePorts(data *common.OrdererData) []string {
	ports := make([]string, 0, 1)
	ports = append(ports, strconv.Itoa(data.Extra.Port)+":"+strconv.Itoa(data.Extra.Port))
	return ports
}

func (dm *DataMachineOrderer) prepareCMDs(data *common.OrdererData) []string {
	cmd := prepareMSPCMDs(data.Service.Name, data.Org, &data.Extra.HFNode)
	genesis := base64.StdEncoding.EncodeToString(data.Genesis)
	genesisCmd := " && echo " + "\"" + genesis + "\"" + " > " + "/var/hyperledger/production/block && " +
		"mkdir -p " + PATHGenesis + " && " +
		"base64 -d /var/hyperledger/production/block > " + PATHGenesis + GenesisFile + " && "
	cmd[2] += genesisCmd
	ordererCmd := " orderer"
	cmd[2] += ordererCmd
	return cmd
}
