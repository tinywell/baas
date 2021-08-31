package docker

import (
	"encoding/json"
	"fmt"
	"strconv"

	"baas/common/tools"
	"baas/internal/model"
	"baas/internal/service/runtime/metadata/common"
	"baas/pkg/runtime"
	"baas/pkg/runtime/docker"
)

// PATH const
const (
	PATHTLS        = "/etc/hyperledger/fabric/tls"
	PATHTLSCert    = PATHTLS + "/server.crt"
	PATHTLSKey     = PATHTLS + "/server.key"
	PATHTLSCA      = PATHTLS + "/ca.crt"
	PATHMSP        = "/etc/hyperledger/fabric/msp"
	PATHOUConfig   = PATHMSP + "/config.yaml"
	PATHDockerDock = "/host/var/run/docker.soc"

	ImagePeer = "hyperledger/fabric-peer"
	PortPeer  = 7051
)

// 默认值
var (
	DefaultEnvPeer = map[string]string{
		"CORE_VM_ENDPOINT":                      "unix:///host/var/run/docker.sock",
		"CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE": "default",
		"CORE_LOGGING_LEVEL":                    "DEBUG",
		"FABRIC_LOGGING_SPEC":                   "DEBUG",
		"CORE_PEER_GOSSIP_USELEADERELECTION":    "true",
		"CORE_PEER_GOSSIP_ORGLEADER":            "false",
		"CORE_PEER_ID":                          "peer.example.com",
		"CORE_PEER_ADDRESS":                     "peer.example.com:7051",
		"CORE_PEER_LISTENADDRESS":               "0.0.0.0:7051",
		"CORE_PEER_CHAINCODEADDRESS":            "peer.example.com:7052",
		"CORE_PEER_CHAINCODELISTENADDRESS":      "0.0.0.0:7052",
		"CORE_PEER_GOSSIP_BOOTSTRAP":            "peer.example.com:7051",
		"CORE_PEER_GOSSIP_EXTERNALENDPOINT":     "peer.example.com:7051",
		"CORE_PEER_LOCALMSPID":                  "PeerMSP",
		"CORE_PEER_TLS_ENABLED":                 "true",
		"CORE_PEER_TLS_CERT_FILE":               PATHTLSCert,
		"CORE_PEER_TLS_KEY_FILE":                PATHTLSKey,
		"CORE_PEER_TLS_ROOTCERT_FILE":           PATHTLSCA,
		"CORE_PEER_MSPCONFIGPATH":               PATHMSP,
	}
)

// CreateData  peer 节点创建数据
func (dm *DataMachinePeer) CreateData(data *common.PeerData) runtime.ServiceMetadata {
	//TODO:
	if data.Service.Runtime != model.RuntimeTypeDocker {
		return nil
	}
	svcData := docker.NewSingleServiceData()
	svcData.Name = data.Service.Name
	if len(data.Extra.Image) == 0 {
		data.Extra.Image = ImagePeer
	}
	if len(data.Extra.Tag) > 0 {
		svcData.Image = data.Extra.Image + ":" + data.Extra.Tag
	} else {
		svcData.Image = data.Extra.Image
	}
	svcData.Ports = preparePorts(data)
	svcData.Volumes = dm.prepareVolumes(data)
	svcData.ENVs = dm.prepareEnvs(data)
	svcData.ExtraHosts = data.ExtraHost
	svcData.Network = prepareNetwork(data.NetworkName, data.Extra.MSPID)
	svcData.CMDs = dm.prepareCMDs(data)

	return svcData
}

// DeleteData ...
func (dm *DataMachinePeer) DeleteData(data *common.PeerData) runtime.ServiceMetadata {
	//TODO:
	return nil
}

func (dm *DataMachinePeer) prepareEnvs(data *common.PeerData) []string {
	envs := tools.CopyStrMap(DefaultEnvPeer)

	envs["CORE_PEER_LOCALMSPID"] = data.Extra.MSPID
	envs["CORE_PEER_ID"] = data.Extra.Name
	envs["CORE_PEER_ADDRESS"] = data.Extra.Endpoint

	envStr := make([]string, 0, len(envs))
	for k, v := range envs {
		envStr = append(envStr, k+"="+v)
	}
	return envStr
}

func (dm *DataMachinePeer) prepareVolumes(data *common.PeerData) []string {
	dc := &model.DataCenterDocker{}
	if err := json.Unmarshal(data.Service.DCMetadata, dc); err != nil {
		//TODO:
		return nil
	}
	vols := []string{}
	if len(dc.Sock) > 0 {
		vols = append(vols, dc.Sock+":"+PATHDockerDock)
	}
	return vols
}

func (dm *DataMachinePeer) prepareCMDs(data *common.PeerData) []string {
	cmd := prepareMSPCMDs(data.Service.Name, data.Org, &data.Extra.HFNode)
	peerCmd := " && peer node start"
	cmd[2] += peerCmd
	return cmd
}

func preparePorts(data *common.PeerData) []string {
	ports := make([]string, 0, 1)
	ports = append(ports, strconv.Itoa(data.Extra.Port)+":"+strconv.Itoa(PortPeer))
	return ports
}

func (dm *DataMachinePeer) prepareNetwork(data *common.PeerData) string {
	return fmt.Sprintf("BaaSNet%04d%s", data.Service.NetworkID, data.Service.MSPID)
}
