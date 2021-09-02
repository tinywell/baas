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
	anchors    map[string][]string
	bootstrap  map[string]string
}

// Builder 运行时相关数据构建
type Builder interface {
	GetDCName(raw []byte) (string, error)
	GetDCData(host request.RuntimeHost) ([]byte, error)
	GetRunner(host request.RuntimeHost) (runtime.ServiceRunner, error)
}
type netService struct {
	builder  Builder
	services fabservices
	genesis  []byte
	runner   map[string]runtime.ServiceRunner // 不同的 datacenter 有不同的 runner 实例（docker 则以 hostname 作为 datacenter ）
}

func (net *netService) Init(req *request.NetInit) error {
	net.builder = CreateBuilder(req.Runtime)
	net.services = fabservices{
		orgs:       make(map[string]*model.FOrganization),
		peers:      make(map[string][]*model.Peer),
		orderers:   make(map[string][]*model.Orderer),
		vmservices: make(map[string]*model.VMService),
		anchors:    make(map[string][]string),
		bootstrap:  make(map[string]string),
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
	err = net.runOrderer(req)
	if err != nil {
		return errors.WithMessage(err, "启动 orderer 节点出错")
	}

	err = net.runPeer(req)
	if err != nil {
		return errors.WithMessage(err, "启动 peer 节点出错")
	}

	return nil
}

func (net *netService) preReq(req *request.NetInit) error {
	if len(req.StateDB) == 0 {
		req.StateDB = model.StateDBLevelDB
	}
	if len(req.CryptoType) == 0 {
		req.CryptoType = model.CryptoTypeSW
	}
	hosts := make(map[string][]byte)

	for _, h := range req.Hosts {
		dcdata, err := net.builder.GetDCData(h)
		if err != nil {
			return err
		}
		hosts[h.Name] = dcdata
		runner, err := net.builder.GetRunner(h)
		if err != nil {
			return err
		}
		net.runner[h.Name] = runner
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
				endpoint := fmt.Sprintf("%s:%d", p.Name, p.Port)
				net.services.peers[r.MSPID][i] = &model.Peer{
					Name:       p.Name,
					DomainName: p.Name,
					Endpoint:   endpoint,
					Port:       p.Port,
					Image:      req.ImagePeer,
					Tag:        req.Version,
				}
				vm := &model.VMService{
					MSPID:      r.MSPID,
					Name:       p.Name,
					Runtime:    model.RuntimeTypeNameValue[req.Runtime],
					LinkType:   model.VMServiceTypePeer,
					DataCenter: p.DataCenter,
					DCMetadata: hosts[p.DataCenter],
				}

				net.services.vmservices[p.Name] = vm
				if p.IsAnchor {
					anc := net.services.anchors[r.MSPID]
					anc = append(anc, endpoint)
					net.services.anchors[r.MSPID] = anc
				}
				if p.IsBootstrap {
					net.services.bootstrap[r.MSPID] = endpoint
				}
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
					Tag:        req.Version,
				}
				vm := &model.VMService{
					MSPID:      r.MSPID,
					Name:       o.Name,
					Runtime:    model.RuntimeTypeNameValue[req.Runtime],
					LinkType:   model.VMServiceTypePeer,
					DataCenter: o.DataCenter,
					DCMetadata: hosts[o.DataCenter],
				}

				net.services.vmservices[o.Name] = vm
			}
		}
	}

	return nil
}

func (net *netService) generateCerts(req *request.NetInit) error {
	msp := newmsp(req.CryptoType) //TODO: nodeou 等配置需要制定并用于生成 msp 实例
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
			// TODO: 此处用于 anchor peers 设置，需要根据请求参数处理
			enpoints := make([]configtx.Endpoint, 0, len(peers))
			for _, p := range peers {
				enpoints = append(enpoints, configtx.Endpoint(p.Endpoint))
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
				OrdererType: req.GenesisConfig.Type,
			}
			name, err := net.builder.GetDCName(ser.DCMetadata)
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
			return errors.WithMessagef(err, "在 %s 启动节点出错", k)
		}
	}

	return nil
}

func (net *netService) runPeer(req *request.NetInit) error {
	ctx := context.Background()
	peers := make(map[string][]*rc.PeerData, 0)

	for k, v := range net.services.peers {
		org := net.services.orgs[k]
		for _, p := range v {
			ser := net.services.vmservices[p.Name]
			boots := net.services.bootstrap[k]
			datap := &rc.PeerData{
				Service:     ser,
				Extra:       p,
				Org:         org,
				NetworkName: req.Network.Name,
				ExtraHost:   []string{}, //TODO: docker 运行时需要准备 extra_hosts
				BootStraps:  boots,
			}
			name, err := net.builder.GetDCName(ser.DCMetadata)
			if err != nil {
				return err
			}
			ops := peers[name]
			ops = append(ops, datap)
			peers[name] = ops
		}
	}

	for k, v := range peers {
		if len(v) == 0 {
			continue
		}
		runner, ok := net.runner[k]
		if !ok {
			return errors.Errorf("runner %s 不存在，无法启动节点", k)
		}
		sr := rs.NewService(model.RuntimeTypeNameValue[req.Runtime], runner)
		err := sr.RunPeers(ctx, v)
		if err != nil {
			return errors.WithMessagef(err, "在 %s 启动节点出错", k)
		}
	}
	return nil
}

func getDCNameFunc(runtime string) func([]byte) (string, error) {
	dockerFunc := func(dcraw []byte) (string, error) {
		dc := &model.DataCenterDocker{}
		err := dc.FromBytes(dcraw)
		if err != nil {
			return "", errors.Errorf("反序列化节点 datacenter 数据出错")
		}
		return dc.Name, nil
	}

	switch runtime {
	case model.RuntimeTypeNameDocker:
		return dockerFunc
	default:
		return dockerFunc
	}
}
