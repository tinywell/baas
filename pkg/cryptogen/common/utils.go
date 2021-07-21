package common

import (
	"os"
	"path/filepath"
	"strconv"
)

// GenerateCaDir 生成临时的 ca 证书存储目录
func GenerateCaDir(networkID int64) (cakeystorepath string, tlscakeystory string, err error) {
	tempDir := makeTempdir(networkID)
	rlcakeystorepath := filepath.Join(tempDir, "cakeystory")
	rltlscakeystory := filepath.Join(tempDir, "tlscakeystory")
	folders := []string{
		rlcakeystorepath,
		rltlscakeystory,
	}
	for _, folder := range folders {
		err := os.MkdirAll(folder, 0755)
		if err != nil {
			return cakeystorepath, tlscakeystory, err
		}
	}
	cakeystorepath, err = filepath.Abs(rlcakeystorepath)
	if err != nil {
		return cakeystorepath, tlscakeystory, err
	}
	tlscakeystory, err = filepath.Abs(rltlscakeystory)
	if err != nil {
		return cakeystorepath, tlscakeystory, err
	}

	return cakeystorepath, tlscakeystory, err
}

// GenerateNodeDir 生成临时的节点证书存储目录
func GenerateNodeDir(networkID int64) (string, string, error) {

	tempDir := makeTempdir(networkID)
	rlnodekeystory := filepath.Join(tempDir, "nodekeystory")
	rlnodetlskeystory := filepath.Join(tempDir, "nodetlskeystory")
	folders := []string{
		rlnodekeystory,
		rlnodetlskeystory,
	}
	for _, folder := range folders {
		err := os.MkdirAll(folder, 0755)
		if err != nil {
			return "", "", err
		}
	}
	nodekeystory, err := filepath.Abs(rlnodekeystory)
	if err != nil {
		return "", "", err
	}
	nodetlskeystory, err := filepath.Abs(rlnodetlskeystory)
	if err != nil {
		return "", "", err
	}
	return nodekeystory, nodetlskeystory, err
}

// GenerateMemberDir 生成临时的成员证书存储目录
func GenerateMemberDir(networkID int64, name string) (string, string, error) {
	tempDir := makeTempdir(networkID)
	keystory := filepath.Join(tempDir, "member", name)
	tlskeystory := filepath.Join(tempDir, "membertls", name)
	folders := []string{
		keystory,
		tlskeystory,
	}
	for _, folder := range folders {
		err := os.MkdirAll(folder, 0755)
		if err != nil {
			return "", "", err
		}
	}
	nodekeystory, err := filepath.Abs(keystory)
	if err != nil {
		return "", "", err
	}
	nodetlskeystory, err := filepath.Abs(tlskeystory)
	if err != nil {
		return "", "", err
	}
	return nodekeystory, nodetlskeystory, err
}

func makeTempdir(netid int64) string {
	dir := os.TempDir()
	tempDir := filepath.Join(dir, "baastemp")
	intermediateDir := ""
	intermediateDir = strconv.Itoa(int(netid))
	return filepath.Join(tempDir, intermediateDir)
}

// CleanupBaastemp 清理临时目录
func CleanupBaastemp(netid int64) {
	tmp := makeTempdir(netid)
	os.RemoveAll(tmp)
}

// Cleanup 清理临时目录
func Cleanup(dir string) {
	os.RemoveAll(dir)
}
