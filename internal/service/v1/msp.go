package v1

import (
	"baas/internal/model"
	"baas/pkg/cryptogen"
	"baas/pkg/cryptogen/common"

	"github.com/pkg/errors"
)

// MSPOrg ...
type MSPOrg struct {
	MSPID      string
	BaseDomain string
}

// MSPNode ...
type MSPNode struct {
	MSPID  string
	Name   string
	NodeOU string
}

// NodeOU
const (
	NodeOUPeer    = "peer"
	NodeOUOrderer = "orderer"
	NodeOUAdmin   = "admin"
	NodeOUClient  = "client"
)

// Node 类型，peer、orderer、user
const (
	NodeTypePeer    = "peer"
	NodeTypeOrderer = "orderer"
	NodeTypeUser    = "user"
)

type msp struct {
	gen cryptogen.Generator
}

func newmsp(crypto string) *msp {
	ct := cryptogen.CryptoTypeSW //
	gen := cryptogen.NewCenarator(ct)
	return &msp{
		gen: gen,
	}
}

func (m *msp) genOrg(org *MSPOrg) (*model.HFNode, error) {
	msporg, err := m.gen.GenerateOrgCA(&common.NodeSpec{
		Organization: org.BaseDomain,
		CommonName:   org.BaseDomain,
	})
	if err != nil {
		return nil, errors.WithMessagef(err, "签发组织 %s(%s) 根证书出错", org.MSPID, org.BaseDomain)
	}
	node := &model.HFNode{
		MSPID:   org.MSPID,
		Name:    org.BaseDomain,
		MSPKey:  msporg.MSPCAKey,
		MSPCert: msporg.MSPCACert,
		TLSKey:  msporg.TLSCAKey,
		TLSCert: msporg.TLSCACert,
	}
	return node, nil
}

func (m *msp) genPeer(org *model.HFNode, mem *MSPNode) (*model.HFNode, error) {
	mem.NodeOU = NodeOUPeer
	return m.genMember(org, mem)
}

func (m *msp) genOrderer(org *model.HFNode, mem *MSPNode) (*model.HFNode, error) {
	mem.NodeOU = NodeOUOrderer
	return m.genMember(org, mem)
}

func (m *msp) genUser(org *model.HFNode, mem *MSPNode, admin bool) (*model.HFNode, error) {
	if admin {
		mem.NodeOU = NodeOUAdmin
	} else {
		mem.NodeOU = NodeOUClient
	}
	return m.genMember(org, mem)
}

func (m *msp) genMember(org *model.HFNode, mem *MSPNode) (*model.HFNode, error) {
	ca := &common.Organization{
		Name:      org.Name,
		MSPID:     org.MSPID,
		TLSCACert: org.TLSCert,
		TLSCAKey:  org.TLSKey,
		MSPCACert: org.MSPCert,
		MSPCAKey:  org.MSPKey,
	}
	spec := &common.NodeSpec{
		Organization:       ca.Name,
		CommonName:         mem.Name,
		OrganizationalUnit: mem.NodeOU,
	}
	mspmem, err := m.gen.GenarateMember(spec, ca)
	if err != nil {
		return nil, errors.WithMessagef(err, "签发成员 %s(%s) 根证书出错", mem.MSPID, mem.Name)
	}
	node := &model.HFNode{
		MSPID:   mem.MSPID,
		Name:    mem.Name,
		MSPKey:  mspmem.MSPKey,
		MSPCert: mspmem.MSPCert,
		TLSKey:  mspmem.TLSKey,
		TLSCert: mspmem.TLSCert,
	}
	return node, nil
}
