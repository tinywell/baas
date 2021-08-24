package configtx

import (
	"crypto/x509"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric-config/configtx"
	"github.com/hyperledger/fabric-config/configtx/membership"
	"github.com/hyperledger/fabric-config/configtx/orderer"
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
			Policies: make(map[string]configtx.Policy),
			Orderer: configtx.Orderer{
				OrdererType:  DefOrdererType,
				BatchTimeout: DefBatchTimeout,
				BatchSize: orderer.BatchSize{
					MaxMessageCount:   DefBatchSizeMaxCount,
					AbsoluteMaxBytes:  DefBatchSizeAbsolute,
					PreferredMaxBytes: DefBatchSizePerffered,
				},
				EtcdRaft: orderer.EtcdRaft{
					Options: orderer.EtcdRaftOptions{},
				},
			},
			Application: configtx.Application{},
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
		Organizations: make([]configtx.Organization, 0, len(members)),
	}
	for _, mem := range members {
		org, err := mem.toOrg()
		if err != nil {
			return err
		}
		org.Policies = mem.toPeerPolicy()
		cg.Organizations = append(cg.Organizations, org)
	}
	c.channel.Consortiums = append(c.channel.Consortiums, cg)
	return nil
}

// SetOrderers 设置共识组
func (c *ConfigTx) SetOrderers(ordrers []*Organization) error {
	for _, o := range ordrers {
		co, err := o.toOrg()
		if err != nil {
			return err
		}
		co.OrdererEndpoints = make([]string, 0, len(o.Endpoints))
		for _, e := range o.Endpoints {
			co.OrdererEndpoints = append(co.OrdererEndpoints, (string)(e))
		}
		co.Policies = o.toOrdererPolicy()
		c.channel.Orderer.Organizations = append(c.channel.Orderer.Organizations, co)
	}
	return nil
}

// SetOrdererConfig 设置共识配置
func (c *ConfigTx) SetOrdererConfig(cfg *OrdererConfig) error {
	if len(cfg.OdererType) > 0 {
		c.channel.Orderer.OrdererType = cfg.OdererType
	}
	if cfg.Cutter.BatchTimeout > time.Duration(0) {
		c.channel.Orderer.BatchTimeout = cfg.Cutter.BatchTimeout
	}
	if cfg.Cutter.BatchSizeAbsolute > 0 {
		c.channel.Orderer.BatchSize.AbsoluteMaxBytes = cfg.Cutter.BatchSizeAbsolute
	}
	if cfg.Cutter.BatchSizePrefer > 0 {
		c.channel.Orderer.BatchSize.PreferredMaxBytes = cfg.Cutter.BatchSizePrefer
	}
	if cfg.Cutter.BatchSizeMaxCount > 0 {
		c.channel.Orderer.BatchSize.MaxMessageCount = cfg.Cutter.BatchSizeMaxCount
	}
	return nil
}

// SetConsenters 设置 raft 节点配置
func (c *ConfigTx) SetConsenters(consenters []*RaftConsentor) error {
	for _, cs := range consenters {
		adds := strings.Split(cs.Address, ":")
		if len(adds) != 2 {
			return errors.Errorf("地址信息格式错误: %s", cs.Address)
		}
		p, err := strconv.Atoi(adds[1])
		if err != nil {
			return errors.WithMessagef(err, "端口信息格式转换出错: %s", adds[1])
		}
		cliCert, err := transCert(cs.ClientCert)
		if err != nil {
			return errors.WithMessage(err, "转换节点 Client 证书格式出错")
		}
		serCert, err := transCert(cs.ServerCert)
		if err != nil {
			return errors.WithMessage(err, "转换节点 Server 证书格式出错")
		}
		rc := orderer.Consenter{
			Address: orderer.EtcdAddress{
				Host: adds[0],
				Port: p,
			},
			ClientTLSCert: cliCert,
			ServerTLSCert: serCert,
		}
		c.channel.Orderer.EtcdRaft.Consenters = append(c.channel.Orderer.EtcdRaft.Consenters, rc)
	}
	return nil
}

// SetApplication 设置应用组
func (c *ConfigTx) SetApplication(members []*Organization) error {
	for _, o := range members {
		co, err := o.toOrg()
		if err != nil {
			return err
		}
		co.AnchorPeers = make([]configtx.Address, 0, len(o.Endpoints))
		for _, e := range o.Endpoints {
			ap, err := e.toAddress()
			if err != nil {
				return errors.WithMessagef(err, "转化组织 %s 的 anchorpeer 地址失败", o.MSPID)
			}
			co.AnchorPeers = append(co.AnchorPeers, ap)
		}
		co.Policies = o.toPeerPolicy()
		c.channel.Application.Organizations = append(c.channel.Application.Organizations, co)
	}
	return nil
}

//============================ Organization ===========================
func (o *Organization) toOrg() (configtx.Organization, error) {
	org := configtx.Organization{
		Name: o.MSPID,
	}
	MSP, err := o.toMSP()
	if err != nil {
		return org, errors.WithMessagef(err, "创建组织（%s） msp 配置出错", o.MSPID)
	}
	org.MSP = MSP
	org.ModPolicy = ModPolicyAdmin
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

func (o *Organization) toOrdererPolicy() map[string]configtx.Policy {
	return map[string]configtx.Policy{
		configtx.ReadersPolicyKey: getOrgPolicy(o.MSPID, []string{NodeOUMember}),
		configtx.WritersPolicyKey: getOrgPolicy(o.MSPID, []string{NodeOUMember}),
		configtx.AdminsPolicyKey:  getOrgPolicy(o.MSPID, []string{NodeOUAdmin}),
	}
}

func (o *Organization) toPeerPolicy() map[string]configtx.Policy {
	return map[string]configtx.Policy{
		configtx.ReadersPolicyKey: getOrgPolicy(o.MSPID, []string{NodeOUAdmin, NodeOUPeer, NodeOUClient}),
		configtx.WritersPolicyKey: getOrgPolicy(o.MSPID, []string{NodeOUAdmin, NodeOUClient}),
		configtx.AdminsPolicyKey:  getOrgPolicy(o.MSPID, []string{NodeOUAdmin}),
	}
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
