package docker

import (
	"baas/internal/model"
	"fmt"
)

func prepareNetwork(network, mspid string) string {
	return fmt.Sprintf("BaaSNet%s%s", network, mspid)
}

func prepareMSPCMDs(service string, org *model.FOrganization, msp *model.HFNode) []string {
	commonCmd := make([]string, 0, 3)
	commonCmd = append(commonCmd, "/bin/sh")
	commonCmd = append(commonCmd, "-c")

	caFilename := "ca." + org.Domian + "-cert.pem" // 需要跟定义 OU 的 config.yaml 中保持一致
	certFilename := service + "-cert.pem"
	keyFilename := service + "_sk"
	tlsCAFilename := "tlsca-" + service + "-cert.pem"
	admincertsName := "Admin@" + org.Domian + "-cert.pem"
	cmd :=
		` pwd && ls && ` +
			"mkdir -p /var/hyperledger/production && " +
			"mkdir -p " + PATHTLS + " && " +
			"echo \"" + msp.TLSCert + "\" > " + PATHTLSCert + " && " +
			"echo \"" + msp.TLSKey + "\" > " + PATHTLSKey + " && " +
			"echo \"" + org.TLSCACert + "\" > " + PATHTLSCA + " && " +
			"rm -rf " + PATHMSP + "/*" + " && " +
			"echo \"" + msp.OUConfig + "\" > " + PATHOUConfig + " && " +
			"mkdir -p " + PATHMSP + "/cacerts" + " && " +
			"rm -rf " + PATHMSP + "/cacerts/*" + " && " +
			"echo \"" + org.CACert + "\" > " + PATHMSP + "/cacerts/" + caFilename + " && " +
			"mkdir -p " + PATHMSP + "/admincerts &&" +
			"rm -rf " + PATHMSP + "/admincerts/*" + " && " +
			"echo \"" + org.AdminCert + "\" > " + PATHMSP + "/admincerts/" + admincertsName + " &&" +
			"mkdir -p " + PATHMSP + "/signcerts" + " && " +
			"rm -rf " + PATHMSP + "/signcerts/*" + " && " +
			"echo \"" + msp.MSPCert + "\" > " + PATHMSP + "/signcerts/" + certFilename + " && " +
			"mkdir -p " + PATHMSP + "/keystore" + " && " +
			"rm -rf " + PATHMSP + "/keystore/*" + " && " +
			"echo \"" + msp.MSPKey + "\" > " + PATHMSP + "/keystore/" + keyFilename + " && " +
			"mkdir -p " + PATHMSP + "/tlscacerts" + " && " +
			"rm -rf " + PATHMSP + "/tlscacerts/*" + " && " +
			"echo \"" + org.TLSCACert + "\" > " + PATHMSP + "/tlscacerts/" + tlsCAFilename +
			// "cat " + PATHOUConfig + " && " +
			""

	commonCmd = append(commonCmd, cmd)
	return commonCmd
}
