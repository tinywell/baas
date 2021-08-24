package v1

import (
	"fmt"
	"time"

	"baas/common"
	"baas/internal/model"
	"baas/internal/model/request"
	"baas/pkg/configtx"

	"github.com/pkg/errors"
)

// Network ...
var Network = netService{}

type fabservices struct {
	orgs     map[string]*model.FOrganization
	peers    map[string][]*model.Peer
	orderers map[string][]*model.Orderer
}
type netService struct {
	services fabservices
	genesis  []byte
}

func (net *netService) Init(req *request.NetInit) error {
	net.services = fabservices{
		orgs:     make(map[string]*model.FOrganization),
		peers:    make(map[string][]*model.Peer),
		orderers: make(map[string][]*model.Orderer),
	}
	err := net.preReq(req)
	if err != nil {
		return errors.WithMessage(err, "请求参数预处理出错")
	}
	// - 签发证书 ：组织根证书 & 组织管理员证书 & 节点证书
	err = net.generateCerts(req)
	if err != nil {
		return errors.WithMessage(err, "签发证书出错")
	}
	// - 准备创世块
	err = net.generateGenesis(req)
	if err != nil {
		return errors.WithMessage(err, "生成创世块出错")
	}
	// - 准备服务创建及启动

	return nil
}

func (net *netService) preReq(req *request.NetInit) error {
	if len(req.StateDB) == 0 {
		req.StateDB = model.StateDBLevelDB
	}
	if len(req.CryptoType) == 0 {
		req.CryptoType = model.CryptoTypeSW
	}
	return nil
}

func (net *netService) generateCerts(req *request.NetInit) error {
	msp := newmsp(req.CryptoType) //TODO:
	// org
	for _, r := range req.Members {
		org, err := msp.genOrg(&MSPOrg{
			MSPID:      r.MSPID,
			BaseDomain: r.Domain,
		})
		if err != nil {
			return err
		}
		forg := &model.FOrganization{
			Name:      r.Domain,
			MSPID:     r.MSPID,
			CACert:    org.MSPCert,
			CAKey:     org.MSPKey,
			TLSCACert: org.TLSCert,
			TLSCAKey:  org.TLSKey,
			Domian:    r.Domain,
		}
		admin, err := msp.genUser(org, &MSPNode{
			MSPID: r.MSPID,
			Name:  fmt.Sprintf("Admin@%s", r.Domain),
		}, true)
		if err != nil {
			return err
		}
		forg.AdminCert = admin.MSPCert
		net.services.orgs[r.MSPID] = forg
		if peers, ok := req.NodePeers[r.MSPID]; ok {
			net.services.peers[r.MSPID] = make([]*model.Peer, len(peers), len(peers))
			for i, p := range peers {
				np, err := msp.genPeer(org, &MSPNode{
					MSPID: r.MSPID,
					Name:  p.Name,
				})
				if err != nil {
					return err
				}
				net.services.peers[r.MSPID][i] = &model.Peer{
					HFNode:     *np,
					Name:       p.Name,
					DomainName: p.Name,
					Endpoint:   fmt.Sprintf("%s:%d", p.Name, p.Port),
					Port:       p.Port,
					Image:      req.ImagePeer,
				}
			}
		}
		if orderers, ok := req.NodeOrderers[r.MSPID]; ok {
			net.services.orderers[r.MSPID] = make([]*model.Orderer, len(orderers), len(orderers))
			for i, o := range orderers {
				no, err := msp.genOrderer(org, &MSPNode{
					MSPID: r.MSPID,
					Name:  o.Name,
				})
				if err != nil {
					return err
				}
				net.services.orderers[r.MSPID][i] = &model.Orderer{
					HFNode:     *no,
					Name:       o.Name,
					DomainName: o.Name,
					Port:       o.Port,
					Endpoint:   fmt.Sprintf("%s:%d", o.Name, o.Port),
					Image:      req.ImagesOrderer,
				}
			}
		}
	}
	return nil
}

func (net *netService) generateGenesis(req *request.NetInit) error {
	config := configtx.NewConfigtx()

	ordererOrgs := make([]*configtx.Organization, 0, len(net.services.orgs))
	for k, v := range net.services.orgs {
		o := configtx.Organization{
			MSPID:     v.MSPID,
			MSPCert:   v.CACert,
			TLSCert:   v.TLSCACert,
			AdminCert: v.AdminCert,
		}

		if orderers, ok := net.services.orderers[k]; ok {
			enpoints := make([]configtx.Endpoint, 0, len(orderers))
			for _, o := range orderers {
				enpoints = append(enpoints, configtx.Endpoint(o.Endpoint))
			}
			o.Endpoints = enpoints
		}
		ordererOrgs = append(ordererOrgs, &o)
	}
	appOrgs := make([]*configtx.Organization, 0, len(net.services.orgs))
	for k, v := range net.services.orgs {
		o := configtx.Organization{
			MSPID:     v.MSPID,
			MSPCert:   v.CACert,
			TLSCert:   v.TLSCACert,
			AdminCert: v.AdminCert,
		}

		if peers, ok := net.services.peers[k]; ok {
			enpoints := make([]configtx.Endpoint, 0, len(peers))
			for _, o := range peers {
				enpoints = append(enpoints, configtx.Endpoint(o.Endpoint))
			}
			o.Endpoints = enpoints
		}
		appOrgs = append(appOrgs, &o)
	}
	err := config.AddConsortium(appOrgs, "")
	if err != nil {
		return errors.WithMessage(err, "添加联盟组织出错")
	}
	err = config.SetOrderers(ordererOrgs)
	if err != nil {
		return errors.WithMessage(err, "添加共识组织出错")
	}
	cfg := configtx.OrdererConfig{
		OdererType: req.GenesisConfig.Type,
		Cutter: configtx.CutterConfig{
			BatchTimeout:      time.Duration(req.GenesisConfig.BatchTimeout),
			BatchSizeAbsolute: uint32(req.GenesisConfig.AbsoluteMaxBytes),
			BatchSizePrefer:   uint32(req.GenesisConfig.PreferredMaxBytes),
			BatchSizeMaxCount: uint32(req.GenesisConfig.MaxMessageCount),
		},
	}
	switch req.GenesisConfig.Type {
	case model.OrdererTypeRaft:
		consenters := make([]*configtx.RaftConsentor, 0)
		for k := range net.services.orgs {
			if orderers, ok := net.services.orderers[k]; ok {
				for _, o := range orderers {
					con := configtx.RaftConsentor{
						Address:    o.Endpoint,
						ServerCert: o.TLSCert,
						ClientCert: o.TLSCert,
					}
					consenters = append(consenters, &con)
				}
			}
		}
		cfg.Raft = consenters
	}

	err = config.SetOrdererConfig(&cfg)
	if err != nil {
		return errors.WithMessage(err, "设置共识配置参数出错")
	}
	genesis, err := config.GenesisBlock(common.SYSChannelName)
	if err != nil {
		return err
	}
	net.genesis = genesis
	return nil
}
