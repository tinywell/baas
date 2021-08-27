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
	CreateData(data *common.PeerData) runtime.ServiceMetadata
	DeleteData(data *common.PeerData) runtime.ServiceMetadata
	// ...
}

// GetPeerWorker ...
func GetPeerWorker(runtime int) PeerDataWorker {
	switch runtime {
	case module.RuntimeTypeDocker:
		return &docker.DataMachinePeer{}
	case module.RuntimeTypeHelm3:
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
	case module.RuntimeTypeDocker:
		return &docker.DataMachineOrderer{}
	// case module.RuntimeTypeHelm3:
	// return &helm3.DataMachinePeer{}
	default:
		return &docker.DataMachineOrderer{}
	}
}
