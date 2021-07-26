package v1

import (
	"context"
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"github.com/tinywell/baas/common/log"
	"github.com/tinywell/baas/internal/service/runtime/metadata"
	"github.com/tinywell/baas/internal/service/runtime/metadata/common"
	"github.com/tinywell/baas/pkg/runtime"
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
	rrC := make(chan *RunningResult, len(peers))
	worker := metadata.GetPeerWorker(s.runtimeType)
	wg := sync.WaitGroup{}
	wg.Add(len(peers))
	for _, pd := range peers {
		go func(pd *common.PeerData) {
			data := worker.PeerCreateData(pd)
			rr := &RunningResult{DataID: pd.Service.Name}
			err := s.runner.Run(ctx, data)
			if err != nil {
				rr.Err = errors.WithMessagef(err, "启动 peer=%s 节点失败", pd.Service.Name)
			} else {
				rr.Msg = "启动 peer=" + pd.Service.Name + " 节点成功"
			}
			rrC <- rr
		}(pd)
	}
	rrs := make([]*RunningResult, 0, len(peers))
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
	fmt.Println(msg)
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
