package service

import (
	"baas/pkg/runtime"
	"baas/pkg/runtime/docker"
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
