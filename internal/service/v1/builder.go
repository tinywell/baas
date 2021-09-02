package v1

import (
	"baas/internal/model"
	"baas/internal/model/request"
	rs "baas/internal/service/runtime/service"
	"baas/pkg/runtime"

	"github.com/pkg/errors"
)

type dockerBuilder struct {
	req request.NetInit
}

func (b *dockerBuilder) GetDCName(raw []byte) (string, error) {
	dc := &model.DataCenterDocker{}
	err := dc.FromBytes(raw)
	if err != nil {
		return "", errors.Errorf("反序列化节点 datacenter 数据出错")
	}
	return dc.Name, nil
}

func (b *dockerBuilder) GetDCData(host request.RuntimeHost) ([]byte, error) {
	dcdocker := model.DataCenterDocker{
		Name: host.Name,
		Host: host.Host,
	}
	raw, err := dcdocker.ToBytes()
	if err != nil {
		return nil, errors.WithMessagef(err, "数据中心 %s 配置数据序列化失败", host.Name)
	}
	return raw, nil
}

func (b *dockerBuilder) GetRunner(host request.RuntimeHost) (runtime.ServiceRunner, error) {
	dcfg := rs.DockerConfig{}
	if len(host.Host) > 0 {
		dcfg.Host = host.Host
	}
	//TODO: docker tls 证书配置(如有)
	docrunner, err := rs.CreateDockerRunner(dcfg)
	if err != nil {
		return nil, errors.Errorf("创建 docker 运行时出错，hostname=%s", host.Name)
	}
	return docrunner, nil
}

type helm3Builder struct {
	req request.NetInit
}

func (b *helm3Builder) GetDCName(raw []byte) (string, error) {
	dc := &model.DataCenterHelm3{}
	err := dc.FromBytes(raw)
	if err != nil {
		return "", errors.Errorf("反序列化节点 datacenter 数据出错")
	}
	return dc.Name, nil
}

func (b *helm3Builder) GetDCData(host request.RuntimeHost) ([]byte, error) {
	dchelm3 := model.DataCenterHelm3{
		Name:     host.Name,
		Repo:     host.HelmConfig.RepoConfig.Repo,
		Kubefile: host.HelmConfig.Kubefile,
	}
	raw, err := dchelm3.ToBytes()
	if err != nil {
		return nil, errors.WithMessagef(err, "数据中心 %s 配置数据序列化失败", host.Name)
	}
	return raw, nil
}

func (b *helm3Builder) GetRunner(host request.RuntimeHost) (runtime.ServiceRunner, error) {
	cfg := rs.Helm3Config{
		Repo:     host.HelmConfig.RepoConfig.Repo,
		Kubefile: host.HelmConfig.Kubefile,
	}
	runner, err := rs.CreateHelm3Runner(cfg)
	if err != nil {
		return nil, errors.Errorf("创建 helm3 运行时出错，hostname=%s", host.Name)
	}
	return runner, nil
}

// CreateBuilder 根据运行时类型创建 Builder 实例
func CreateBuilder(runtime string) Builder {
	switch runtime {
	case model.RuntimeTypeNameDocker:
		return &dockerBuilder{}
	case model.RuntimeTypeNameHelm3:
		return &helm3Builder{}
	default:
		return &dockerBuilder{}
	}
}
