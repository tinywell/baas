package common

import module "github.com/tinywell/baas/internal/model"

// PeerData ...
type PeerData struct {
	Service     *module.VMService
	Extra       *module.Peer
	Org         *module.FOrganization
	NetworkName string
	ExtraHost   []string
	BootStraps  []string
}
