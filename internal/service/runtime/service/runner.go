package service

import (
	"baas/pkg/runtime"
	"baas/pkg/runtime/docker"
	"baas/pkg/runtime/helm3"
)

// DockerConfig ..
type DockerConfig struct {
	Host        string
	TLS         bool
	TLSCertPath string
	TLSKeyPath  string
	TLSCAPath   string
}

// CreateDockerRunner ..
func CreateDockerRunner(cfg DockerConfig) (runtime.ServiceRunner, error) {
	opts := []docker.Option{}

	if len(cfg.Host) > 0 {
		opt := docker.WithConfig(cfg.Host, "")
		opts = append(opts, opt)
	}
	if cfg.TLS {
		opt := docker.WithTLS(cfg.TLSCAPath, cfg.TLSCertPath, cfg.TLSKeyPath)
		opts = append(opts, opt)
	}
	cli, err := docker.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	return cli, nil
}

// Helm3Config helm3 运行时初始化配置
type Helm3Config struct {
	Repo     string
	Kubefile string
}

// CreateHelm3Runner 创建 helm3 运行时
func CreateHelm3Runner(cfg Helm3Config) (runtime.ServiceRunner, error) {
	opts := []helm3.Option{}
	opts = append(opts, helm3.WithKubeConfig(cfg.Kubefile, ""))
	opts = append(opts, helm3.WithRepoConfig(helm3.RepoConfig{
		RepoURL: cfg.Repo,
	}))
	cli, err := helm3.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	return cli, nil
}
