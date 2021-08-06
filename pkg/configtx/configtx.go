package configtx

import (
	"crypto/x509"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric-config/configtx"
	"github.com/hyperledger/fabric-config/configtx/membership"
	"github.com/pkg/errors"
)

// ConfigTx ..
type ConfigTx struct {
	channel *configtx.Channel
}

// NewConfigtx 。。
func NewConfigtx() *ConfigTx {
	return &ConfigTx{
		channel: &configtx.Channel{
			Consortiums:  make([]configtx.Consortium, 0),
			Capabilities: make([]string, 0),
			Policies:     make(map[string]configtx.Policy),
			Orderer:      configtx.Orderer{},
			Application:  configtx.Application{},
		},
	}
}

// AddConsortium 增加联盟组织信息
func (c *ConfigTx) AddConsortium(members []*Organization, name string) error {
	if len(name) == 0 {
		name = ConsortiumName
	}
	cg := configtx.Consortium{
		Name:          name,
		Organizations: make([]configtx.Organization, 0),
	}
	for _, mem := range members {
		org, err := mem.toOrg()
		if err != nil {
			return err
		}
		cg.Organizations = append(cg.Organizations, org)
	}
	c.channel.Consortiums = append(c.channel.Consortiums, cg)
	return nil
}

// SetOrderers 设置共识组
func (c *ConfigTx) SetOrderers(ordrers []*Organization) error {
	return nil
}

// SetOrdererConfig 设置共识配置
func (c *ConfigTx) SetOrdererConfig(cfg *OrdererConfig) error {
	return nil
}

// SetApplication 设置应用组
func (c *ConfigTx) SetApplication(members []*Organization) error {
	return nil
}

//============================ Organization ===========================
func (o *Organization) toOrg() (configtx.Organization, error) {
	org := configtx.Organization{
		Name: o.MSPID,
	}
	MSP, err := o.toMSP()
	if err != nil {
		return org, errors.WithMessage(err, "创建组织 msp 配置出错")
	}
	org.MSP = MSP
	org.ModPolicy = ModPolicy
	return org, nil
}

func (o *Organization) toMSP() (configtx.MSP, error) {
	msp := configtx.MSP{
		Name: o.MSPID,
	}
	rootc, err := transCert(o.MSPCert)
	if err != nil {
		return msp, errors.WithMessagef(err, "解析组织(%s)根证书出错", o.MSPID)
	}
	msp.RootCerts = []*x509.Certificate{rootc}
	adminc, err := transCert(o.AdminCert)
	if err != nil {
		return msp, errors.WithMessagef(err, "解析组织(%s)管理员证书出错", o.MSPID)
	}
	msp.Admins = []*x509.Certificate{adminc}
	tlsrootc, err := transCert(o.TLSCert)
	if err != nil {
		return msp, errors.WithMessagef(err, "解析组织(%s) TLS 根证书出错", o.MSPID)
	}
	msp.TLSRootCerts = []*x509.Certificate{tlsrootc}
	if o.NodeOU {
		msp.NodeOUs = nodeOU(rootc)
	}
	msp.CryptoConfig = membership.CryptoConfig{
		SignatureHashFamily:            HashFamily,
		IdentityIdentifierHashFunction: HashFunc,
	}
	return msp, nil
}

func nodeOU(cert *x509.Certificate) membership.NodeOUs {
	return membership.NodeOUs{
		Enable: true,
		ClientOUIdentifier: membership.OUIdentifier{
			Certificate:                  cert,
			OrganizationalUnitIdentifier: "client",
		},
		PeerOUIdentifier: membership.OUIdentifier{
			Certificate:                  cert,
			OrganizationalUnitIdentifier: "peer",
		},
		AdminOUIdentifier: membership.OUIdentifier{
			Certificate:                  cert,
			OrganizationalUnitIdentifier: "admin",
		},
		OrdererOUIdentifier: membership.OUIdentifier{
			Certificate:                  cert,
			OrganizationalUnitIdentifier: "orderer",
		},
	}
}

func (e Endpoint) toAddress() (configtx.Address, error) {
	addr := configtx.Address{}
	ss := strings.Split(string(e), ":")
	if len(ss) != 2 {
		return addr, errors.Errorf("地址 %s 不是 host:port 格式", e)
	}
	port, err := strconv.Atoi(ss[1])
	if err != nil {
		return addr, errors.WithMessagef(err, "端口参数 %s 不是数值", ss[1])
	}
	addr.Host = ss[0]
	addr.Port = port
	return addr, nil
}
