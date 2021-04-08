package common

import "github.com/tinywell/baas/internal/module"

// PeerData ...
type PeerData struct {
	Service     *module.VMService
	Extra       *module.Peer
	Org         *module.FOrganization
	NetworkName string
	ExtraHost   []string
}
