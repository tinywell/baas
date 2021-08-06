package cryptogen

import (
	"baas/pkg/cryptogen/common"
	"baas/pkg/cryptogen/sw"
)

// CryptoType crypto 类型
type CryptoType string

// crypto 类型
const (
	CryptoTypeSW CryptoType = "sw"
	CryptoTypeGM CryptoType = "gm"
)

// Generator 证书颁发器
type Generator interface {
	GenerateOrgCA(org *common.NodeSpec) (common.Organization, error)                        // 签发组织根证书
	GenarateMember(member *common.NodeSpec, CA *common.Organization) (common.Member, error) // 签发组织成员证书
}

type options struct {
	EnableNodeOU bool
}

// Option ...
type Option func(*options)

// NewCenarator ...
func NewCenarator(ct CryptoType, opts ...Option) Generator {
	// TODO:
	switch ct {
	case CryptoTypeSW:
		return &sw.Gen{}
	case CryptoTypeGM:
		return nil
	default:
		return &sw.Gen{}
	}
}

// WithNodeOU 启用 nodeou
func WithNodeOU() Option {
	return func(o *options) {
		o.EnableNodeOU = true
	}
}
