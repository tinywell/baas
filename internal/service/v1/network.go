package v1

import (
	"context"
	"fmt"
	"time"

	"baas/common"
	"baas/internal/model"
	"baas/internal/model/request"
	"baas/pkg/configtx"
	"baas/pkg/runtime"

	rc "baas/internal/service/runtime/metadata/common"
	rs "baas/internal/service/runtime/service"

	"github.com/pkg/errors"
)

// Network ...
var Network = netService{}

type fabservices struct {
	orgs       map[string]*model.FOrganization
	peers      map[string][]*model.Peer
	orderers   map[string][]*model.Orderer
	vmservices map[string]*model.VMService
}
type netService struct {
	services fabservices
	genesis  []byte
	runner   map[string]runtime.ServiceRunner // 不同的 datacenter 有不同的 runner 实例（docker 则以 hostname 作为 datacenter ）
}

func (net *netService) Init(req *request.NetInit) error {
	net.services = fabservices{
		orgs:       make(map[string]*model.FOrganization),
		peers:      make(map[string][]*model.Peer),
		orderers:   make(map[string][]*model.Orderer),
		vmservices: make(map[string]*model.VMService),
	}
	net.runner = make(map[string]runtime.ServiceRunner)

	err := net.preReq(req)
	if err != nil {
		return errors.WithMessage(err, "请求参数预处理出错")
	}
	// - 签发证书 ：组织根证书 & 组织管理员证书 & 节点证书
	err = net.generateCerts(req)
	if err != nil {
		return errors.WithMessage(err, "签发证书出错")
	}

	// - 网络及节点数据入库

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
	hosts := make(map[int]*model.DataCenterDocker)
	if req.Runtime == model.RuntimeTypeNameDocker {
		for id, h := range req.Hosts {
			dcdocker := model.DataCenterDocker{
				Name: h.Hostname,
				Host: h.IP,
				Port: int(h.Port),
			}
			hosts[id] = &dcdocker
			dcfg := rs.DockerConfig{}
			if len(h.IP) > 0 {
				dcfg.Host = fmt.Sprintf("%s:%d", h.IP, h.Port)
			}
			//TODO: docker tls 证书配置
			runner, err := rs.CreateDockerRunner(dcfg)
			if err != nil {
				return errors.Errorf("创建 docker 运行时出错，hostname=%s", h.Hostname)
			}
			net.runner[h.Hostname] = runner
		}
	}
	for _, r := range req.Members {
		// 组织
		forg := &model.FOrganization{
			Name:   r.Domain,
			MSPID:  r.MSPID,
			Domian: r.Domain,
		}
		net.services.orgs[r.MSPID] = forg

		// peer 节点
		if peers, ok := req.NodePeers[r.MSPID]; ok {
			net.services.peers[r.MSPID] = make([]*model.Peer, len(peers), len(peers))
			for i, p := range peers {
				net.services.peers[r.MSPID][i] = &model.Peer{
					Name:       p.Name,
					DomainName: p.Name,
					Endpoint:   fmt.Sprintf("%s:%d", p.Name, p.Port),
					Port:       p.Port,
					Image:      req.ImagePeer,
				}
				vm := &model.VMService{
					MSPID:      r.MSPID,
					Name:       p.Name,
					Runtime:    model.RuntimeTypeNameValue[req.Runtime],
					LinkType:   model.VMServiceTypePeer,
					DataCenter: p.DataCenter,
					DCID:       p.HostID,
				}
				switch req.Runtime {
				case model.RuntimeTypeNameDocker:
					h := hosts[int(p.HostID)]
					raw, err := h.ToBytes()
					if err != nil {
						return errors.WithMessagef(err, "序列化主机数据出错，主机 = %s", h.Name)
					}
					vm.DCMetadata = raw
				case model.RuntimeTypeNameKubenetes:
				}
				net.services.vmservices[p.Name] = vm
			}
		}

		// orderer 节点
		if orderers, ok := req.NodeOrderers[r.MSPID]; ok {
			net.services.orderers[r.MSPID] = make([]*model.Orderer, len(orderers), len(orderers))
			for i, o := range orderers {
				net.services.orderers[r.MSPID][i] = &model.Orderer{
					Name:       o.Name,
					DomainName: o.Name,
					Port:       o.Port,
					Endpoint:   fmt.Sprintf("%s:%d", o.Name, o.Port),
					Image:      req.ImagesOrderer,
				}
				vm := &model.VMService{
					MSPID:      r.MSPID,
					Name:       o.Name,
					Runtime:    model.RuntimeTypeNameValue[req.Runtime],
					LinkType:   model.VMServiceTypePeer,
					DataCenter: o.DataCenter,
					DCID:       o.HostID,
				}
				switch req.Runtime {
				case model.RuntimeTypeNameDocker:
					h := hosts[int(o.HostID)]
					raw, err := h.ToBytes()
					if err != nil {
						return errors.WithMessagef(err, "序列化主机数据出错，主机 = %s", h.Name)
					}
					vm.DCMetadata = raw
				case model.RuntimeTypeNameKubenetes:
				}
				net.services.vmservices[o.Name] = vm
			}
		}
	}

	return nil
}

func (net *netService) generateCerts(req *request.NetInit) error {
	msp := newmsp(req.CryptoType) //TODO:
	// org
	for _, r := range net.services.orgs {
		org, err := msp.genOrg(&MSPOrg{
			MSPID:      r.MSPID,
			BaseDomain: r.Name,
		})
		if err != nil {
			return err
		}
		r.CACert = org.MSPCert
		r.CAKey = org.MSPKey
		r.TLSCACert = org.TLSCert
		r.TLSCAKey = org.TLSKey

		admin, err := msp.genUser(org, &MSPNode{
			MSPID: r.MSPID,
			Name:  fmt.Sprintf("Admin@%s", r.Name),
		}, true)
		if err != nil {
			return err
		}
		r.AdminCert = admin.MSPCert

		if peers, ok := net.services.peers[r.MSPID]; ok {
			for _, p := range peers {
				np, err := msp.genPeer(org, &MSPNode{
					MSPID: r.MSPID,
					Name:  p.Name,
				})
				if err != nil {
					return err
				}
				p.HFNode = *np
			}
		}
		if orderers, ok := net.services.orderers[r.MSPID]; ok {
			for _, o := range orderers {
				no, err := msp.genOrderer(org, &MSPNode{
					MSPID: r.MSPID,
					Name:  o.Name,
				})
				if err != nil {
					return err
				}
				o.HFNode = *no
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

func (net *netService) runOrderer(req *request.NetInit) error {
	ctx := context.Background()
	var GetDCName func([]byte) (string, error)

	switch req.Runtime {
	case model.RuntimeTypeNameDocker:
		GetDCName = func(dcraw []byte) (string, error) {
			dc := &model.DataCenterDocker{}
			err := dc.FromBytes(dcraw)
			if err != nil {
				return "", errors.Errorf("反序列化节点 datacenter 数据出错")
			}
			return dc.Name, nil
		}
	}

	orderers := make(map[string][]*rc.OrdererData, 0)
	for k, v := range net.services.orderers {
		org := net.services.orgs[k]
		for _, o := range v {
			ser := net.services.vmservices[o.Name]
			datao := &rc.OrdererData{
				Service:     ser,
				Extra:       o,
				Org:         org,
				NetworkName: req.Network.Name,
				ExtraHost:   []string{}, //TODO: docker 运行时需要准备 extra_hosts
				Genesis:     net.genesis,
			}
			name, err := GetDCName(ser.DCMetadata)
			if err != nil {
				return err
			}
			ods := orderers[name]
			ods = append(ods, datao)
			orderers[name] = ods

		}
	}
	for k, v := range orderers {
		if len(v) == 0 {
			continue
		}
		runner, ok := net.runner[k]
		if !ok {
			return errors.Errorf("runner %s 不存在，无法启动节点", k)
		}
		sr := rs.NewService(model.RuntimeTypeNameValue[req.Runtime], runner)
		err := sr.RunOrderers(ctx, v)
		if err != nil {
			return errors.Errorf("在 %s 启动节点出错", k)
		}
	}

	return nil
}

func (net *netService) runPeer(req *request.NetInit) error {
	peers := make([]*rc.PeerData, 0)
	for k, v := range net.services.peers {
		org := net.services.orgs[k]
		for _, p := range v {
			datap := &rc.PeerData{
				Service:     net.services.vmservices[p.Name],
				Extra:       p,
				Org:         org,
				NetworkName: req.Network.Name,
				ExtraHost:   []string{},           //TODO: docker 运行时需要准备 extra_hosts
				BootStraps:  []string{p.Endpoint}, //TODO: peer 指定 bootstrap
			}
			peers = append(peers, datap)
		}
	}
	return nil
}
