package metadata

import (
	"github.com/tinywell/baas/internal/module"
	"github.com/tinywell/baas/internal/service/v1/metadata/common"
	"github.com/tinywell/baas/internal/service/v1/metadata/docker"
	"github.com/tinywell/baas/pkg/runtime"
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
	default:
		return &docker.DataMachine{}
	}
}
