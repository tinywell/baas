package service

import (
	"context"
	"fmt"
	"sync"

	"baas/common/log"
	"baas/internal/service/runtime/metadata"
	"baas/internal/service/runtime/metadata/common"
	"baas/pkg/runtime"

	"github.com/pkg/errors"
)

// Service ...
type Service struct {
	runner      runtime.ServiceRunner
	runtimeType int
}

// RunningResult ...
type RunningResult struct {
	DataID string
	Err    error
	Msg    string
}

// RunPeers 批量启动 peer 节点
func (s *Service) RunPeers(ctx context.Context, peers []*common.PeerData) error {
	worker := metadata.GetPeerWorker(s.runtimeType)
	services := make([]runtime.ServiceMetadata, 0, len(peers))
	for _, pd := range peers {
		data := worker.PeerCreateData(pd)
		services = append(services, data)
	}
	err := s.runServices(ctx, services)
	if err != nil {
		return errors.WithMessage(err, "启动 peer 节点出错")
	}
	return nil
}

func (s *Service) runServices(ctx context.Context, services []runtime.ServiceMetadata) error {
	rrC := make(chan *RunningResult, len(services))
	wg := sync.WaitGroup{}
	wg.Add(len(services))
	for _, ser := range services {
		go func(ser runtime.ServiceMetadata) {
			rr := &RunningResult{DataID: ser.DataID()}
			err := s.runner.Run(ctx, ser)
			if err != nil {
				rr.Err = errors.WithMessagef(err, "服务 %s 执行 [%s] 失败", ser.DataID(), ser.Action())
			} else {
				rr.Msg = fmt.Sprintf("服务 %s 执行 [%s] 成功", ser.DataID(), ser.Action())
			}
			rrC <- rr
		}(ser)
	}
	rrs := make([]*RunningResult, 0, len(services))
	go func() {
		for rr := range rrC {
			rrs = append(rrs, rr)
			wg.Done()
		}
	}()
	wg.Wait()
	close(rrC)
	msg, err := countResult(rrs)
	if err != nil {
		return err
	}
	log.Info(msg)
	return nil
}

func countResult(rrs []*RunningResult) (string, error) {
	msg := ""
	errs := []error{}
	for i, r := range rrs {
		if r.Err != nil {
			errs = append(errs, r.Err)
		} else {
			msg += fmt.Sprintf("[%d] %s ID=%s", i, r.Msg, r.DataID)
		}
	}
	if len(errs) > 0 {
		errmsg := fmt.Sprintf("发生了 %d 个错误:\n", len(errs))
		for i, e := range errs {
			errmsg += fmt.Sprintf("\t[%d] %s\n", i+1, e.Error())
		}
		return msg, errors.New(errmsg)
	}
	return msg, nil
}
