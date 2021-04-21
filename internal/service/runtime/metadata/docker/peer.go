package docker

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/tinywell/baas/common/tools"
	"github.com/tinywell/baas/internal/module"
	"github.com/tinywell/baas/internal/service/runtime/metadata/common"
	"github.com/tinywell/baas/pkg/runtime"
	"github.com/tinywell/baas/pkg/runtime/docker"
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

// PeerCreateData  peer 节点创建数据
func (dm *DataMachine) PeerCreateData(data *common.PeerData) runtime.ServiceMetadata {
	//TODO:
	if data.Service.Runtime != module.RuntimeTypeDocker {
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
	svcData.Volumes = prepareVolumes(data)
	svcData.ENVs = prepareEnvs(data)
	svcData.ExtraHosts = data.ExtraHost
	svcData.Network = prepareNetwork(data)
	svcData.CMDs = prepareCMDs(data)

	return svcData
}

// PeerDeleteData ...
func (dm *DataMachine) PeerDeleteData(data *common.PeerData) runtime.ServiceMetadata {
	//TODO:
	return nil
}

func prepareEnvs(data *common.PeerData) []string {
	envs := tools.CopyStrMap(DefaultEnvPeer)
	// TODO: 补充/覆盖 env
	envs["CORE_PEER_LOCALMSPID"] = data.Extra.MSPID
	envs["CORE_PEER_ID"] = data.Extra.Name
	envs["CORE_PEER_ADDRESS"] = data.Extra.Endpoint

	envStr := make([]string, 0, len(envs))
	for k, v := range envs {
		envStr = append(envStr, k+"="+v)
	}
	return envStr
}

func prepareVolumes(data *common.PeerData) []string {
	dc := &module.DataCenterDocker{}
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

func prepareCMDs(data *common.PeerData) []string {
	peerCmd := make([]string, 0, 3)
	peerCmd = append(peerCmd, "/bin/sh")
	peerCmd = append(peerCmd, "-c")

	caFilename := "ca." + data.Org.Domian + "-cert.pem" // 需要跟定义 OU 的 config.yaml 中保持一致
	certFilename := data.Service.Name + "-cert.pem"
	keyFilename := data.Service.Name + "_sk"
	tlsCAFilename := "tlsca-" + data.Service.Name + "-cert.pem"
	cmd :=
		` pwd && ls &&` +
			"mkdir -p " + PATHTLS + " && " +
			"echo \"" + data.Extra.TLSCert + "\" > " + PATHTLSCert + " && " +
			"echo \"" + data.Extra.TLSKey + "\" > " + PATHTLSKey + " && " +
			"echo \"" + data.Org.TLSCACert + "\" > " + PATHTLSCA + " && " +
			"rm -rf " + PATHMSP + "/*" + " && " +
			"echo \"" + data.Extra.OUConfig + "\" > " + PATHOUConfig + " && " +
			"mkdir -p " + PATHMSP + "/cacerts" + " && " +
			"rm -rf " + PATHMSP + "/cacerts/*" + " && " +
			"echo \"" + data.Org.CACert + "\" > " + PATHMSP + "/cacerts/" + caFilename + " && " +
			"mkdir -p " + PATHMSP + "/signcerts" + " && " +
			"rm -rf " + PATHMSP + "/signcerts/*" + " && " +
			"echo \"" + data.Extra.MSPCert + "\" > " + PATHMSP + "/signcerts/" + certFilename + " && " +
			"mkdir -p " + PATHMSP + "/keystore" + " && " +
			"rm -rf " + PATHMSP + "/keystore/*" + " && " +
			"echo \"" + data.Extra.MSPKey + "\" > " + PATHMSP + "/keystore/" + keyFilename + " && " +
			"mkdir -p " + PATHMSP + "/tlscacerts" + " && " +
			"rm -rf " + PATHMSP + "/tlscacerts/*" + " && " +
			"echo \"" + data.Org.TLSCACert + "\" > " + PATHMSP + "/tlscacerts/" + tlsCAFilename + " && " +
			"cat " + PATHOUConfig + " && " +
			"peer node start"

	peerCmd = append(peerCmd, cmd)
	return peerCmd
}

func preparePorts(data *common.PeerData) []string {
	ports := make([]string, 0, 1)
	if data.Extra.EXTPort > 0 {
		ports = append(ports, strconv.Itoa(data.Extra.EXTPort)+":"+strconv.Itoa(data.Extra.Port))
	} else {
		ports = append(ports, strconv.Itoa(data.Extra.Port)+":"+strconv.Itoa(data.Extra.Port))
	}
	return ports
}

func prepareNetwork(data *common.PeerData) string {
	return fmt.Sprintf("BaaSNet%04d%s", data.Service.NetworkID, data.Service.MSPID)
}
