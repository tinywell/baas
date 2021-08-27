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
	PATHGenesis = "/var/hyperledger/orderer/orderer.genesis.block"

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
		"ORDERER_GENERAL_GENESISFILE":               PATHGenesis,
		"ORDERER_GENERAL_LOCALMSPID":                "OrdererMSP",
		"ORDERER_GENERAL_LOCALMSPDIR":               "/var/hyperledger/orderer/msp",
		"ORDERER_OPERATIONS_LISTENADDRESS":          "0.0.0.0:17050",
		"ORDERER_GENERAL_TLS_ENABLED":               "true",
		"ORDERER_GENERAL_TLS_PRIVATEKEY":            "/var/hyperledger/orderer/tls/server.key",
		"ORDERER_GENERAL_TLS_CERTIFICATE":           "/var/hyperledger/orderer/tls/server.crt",
		"ORDERER_GENERAL_TLS_ROOTCAS":               "[/var/hyperledger/orderer/tls/ca.crt]",
		"ORDERER_KAFKA_TOPIC_REPLICATIONFACTOR":     "1",
		"ORDERER_KAFKA_VERBOSE":                     "true",
		"ORDERER_GENERAL_CLUSTER_CLIENTCERTIFICATE": "/var/hyperledger/orderer/tls/server.crt",
		"ORDERER_GENERAL_CLUSTER_CLIENTPRIVATEKEY":  "/var/hyperledger/orderer/tls/server.key",
		"ORDERER_GENERAL_CLUSTER_ROOTCAS":           "[/var/hyperledger/orderer/tls/ca.crt]",
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

	envStr := make([]string, 0, len(envs))
	for k, v := range envs {
		envStr = append(envStr, k+"="+v)
	}
	return envStr
}

func (dm *DataMachineOrderer) preparePorts(data *common.OrdererData) []string {
	ports := make([]string, 0, 1)
	ports = append(ports, strconv.Itoa(data.Extra.Port)+":"+strconv.Itoa(PortOrderer))
	return ports
}

func (dm *DataMachineOrderer) prepareCMDs(data *common.OrdererData) []string {
	cmd := prepareMSPCMDs(data.Service.Name, data.Org, &data.Extra.HFNode)
	genesis := base64.StdEncoding.EncodeToString(data.Genesis)
	genesisCmd := " && echo " + "\"" + genesis + "\"" + " > " + "/var/hyperledger/production/block && " +
		"base64 -d /var/hyperledger/production/block > " + PATHGenesis + " &&"
	cmd = append(cmd, genesisCmd)
	peerCmd := " orderer"
	cmd = append(cmd, peerCmd)
	return cmd
}
