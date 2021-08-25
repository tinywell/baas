package common

import module "baas/internal/model"

// PeerData ...
type PeerData struct {
	Service     *module.VMService
	Extra       *module.Peer
	Org         *module.FOrganization
	NetworkName string
	ExtraHost   []string
	BootStraps  []string
}

// OrdererData ...
type OrdererData struct {
	Service     *module.VMService
	Extra       *module.Orderer
	Org         *module.FOrganization
	NetworkName string
	ExtraHost   []string
}
