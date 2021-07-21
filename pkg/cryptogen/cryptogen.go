package cryptogen

import "github.com/tinywell/baas/pkg/cryptogen/common"

// Generator 证书颁发器
type Generator interface {
	GenerateOrgCA(org common.NodeSpec) (common.Organization, error)                       // 签发组织根证书
	GenarateMember(member common.NodeSpec, CA common.Organization) (common.Member, error) // 签发组织成员证书
}
