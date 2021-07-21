package sw

import (
	"testing"

	"github.com/tinywell/baas/pkg/cryptogen/common"
)

func TestGen_GenerateOrgCA(t *testing.T) {
	spec := &common.NodeSpec{
		Organization: "org1.example.com",
		CommonName:   "org1.example.com",
		Country:      "CN",
		Province:     "HuBei",
		Locality:     "WuHan",
	}
	g := &Gen{}
	org, err := g.GenerateOrgCA(spec)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", org)
	memSpec := &common.NodeSpec{
		Organization:       "org1.example.com",
		CommonName:         "peerx.org1.example.com",
		OrganizationalUnit: "peer",
		Country:            "CN",
		Province:           "HuBei",
		Locality:           "WuHan",
	}
	mem, err := g.GenarateMember(memSpec, &org)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", mem)
}
