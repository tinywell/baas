package common

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// GenerateCaDir 生成临时的 ca 证书存储目录
func GenerateCaDir(name string) (cakeystorepath string, tlscakeystory string, err error) {
	tempDir := makeTempdir(name)
	keystore := filepath.Join(tempDir, "cakeystory")
	tlskeystore := filepath.Join(tempDir, "tlscakeystory")
	return generateDir(keystore, tlskeystore)

}

// GenerateMemberDir 生成临时的成员证书存储目录
func GenerateMemberDir(orgname string, memname string) (string, string, error) {
	tempDir := makeTempdir(orgname)
	keystore := filepath.Join(tempDir, "member", memname)
	tlskeystore := filepath.Join(tempDir, "membertls", memname)
	return generateDir(keystore, tlskeystore)
}

func generateDir(keystore, tlskeystore string) (string, string, error) {
	folders := []string{
		keystore,
		tlskeystore,
	}
	for _, folder := range folders {
		err := os.MkdirAll(folder, 0755)
		if err != nil {
			return "", "", err
		}
	}
	nodekeystore, err := filepath.Abs(keystore)
	if err != nil {
		return "", "", err
	}
	nodetlskeystore, err := filepath.Abs(tlskeystore)
	if err != nil {
		return "", "", err
	}
	return nodekeystore, nodetlskeystore, err
}

func makeTempdir(name string) string {
	dir := os.TempDir()
	tempDir := filepath.Join(dir, "baastemp")
	intermediateDir := name
	return filepath.Join(tempDir, intermediateDir)
}

// CleanupBaastemp 清理临时目录
func CleanupBaastemp(name string) {
	tmp := makeTempdir(name)
	os.RemoveAll(tmp)
}

// Cleanup 清理临时目录
func Cleanup(dir string) {
	os.RemoveAll(dir)
}

// LoadPrivateKey 从 filekeystore 中读取私钥文件数据
func LoadPrivateKey(keystore string) ([]byte, error) {
	var rawKey []byte
	var err error

	walkFunc := func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, "_sk") {
			rawKey, err = ioutil.ReadFile(path)
			if err != nil {
				return err
			}
		}
		return nil
	}

	err = filepath.Walk(keystore, walkFunc)
	if err != nil {
		return nil, err
	}
	return rawKey, nil
}
