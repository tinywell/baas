package metadata

import (
	module "baas/internal/model"
	"baas/internal/service/runtime/metadata/common"
	"baas/internal/service/runtime/metadata/docker"
	"baas/internal/service/runtime/metadata/helm3"
	"baas/pkg/runtime"
)

// PeerDataWorker ...
type PeerDataWorker interface {
	PeerCreateData(data *common.PeerData) runtime.ServiceMetadata
	PeerDeleteData(data *common.PeerData) runtime.ServiceMetadata
	// ...
}

// GetPeerWorker ...
func GetPeerWorker(runtime int) PeerDataWorker {
	switch runtime {
	case module.RuntimeTypeDocker:
		return &docker.DataMachine{}
	case module.RuntimeTypeHelm3:
		return &helm3.DataMachine{}
	default:
		return &docker.DataMachine{}
	}
}
