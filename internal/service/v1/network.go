package v1

import "github.com/tinywell/baas/internal/model/request"

// Network ...
var Network = netService{}

type netService struct{}

func (net *netService) Init(req *request.NetInit) error {
	return nil
}
