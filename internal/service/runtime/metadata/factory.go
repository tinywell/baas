package metadata

import (
	"baas/internal/model"
	"baas/internal/service/runtime/metadata/common"
	"baas/internal/service/runtime/metadata/docker"
	"baas/internal/service/runtime/metadata/helm3"
	"baas/pkg/runtime"
)

// PeerDataWorker ...
type PeerDataWorker interface {
	CreateData(data *common.PeerData) runtime.ServiceMetadata
	DeleteData(data *common.PeerData) runtime.ServiceMetadata
	// ...
}

// GetPeerWorker ...
func GetPeerWorker(runtime int) PeerDataWorker {
	switch runtime {
	case model.RuntimeTypeDocker:
		return &docker.DataMachinePeer{}
	case model.RuntimeTypeHelm3:
		return &helm3.DataMachinePeer{}
	default:
		return &docker.DataMachinePeer{}
	}
}

// OrdererDataWorker ...
type OrdererDataWorker interface {
	CreateData(data *common.OrdererData) runtime.ServiceMetadata
	DeleteData(data *common.OrdererData) runtime.ServiceMetadata
	// ...
}

// GetOrdererWorker ...
func GetOrdererWorker(runtime int) OrdererDataWorker {
	switch runtime {
	case model.RuntimeTypeDocker:
		return &docker.DataMachineOrderer{}
	case model.RuntimeTypeHelm3:
		return &helm3.DataMachineOrderer{}
	default:
		return &docker.DataMachineOrderer{}
	}
}
