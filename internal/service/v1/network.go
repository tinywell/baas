package v1

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/tinywell/baas/internal/model"
	"github.com/tinywell/baas/internal/model/request"
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
