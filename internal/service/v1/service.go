package v1

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
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

// RunPeer 启动 peer 节点
func (s *Service) RunPeer(ctx context.Context, peers []*common.PeerData) error {
	rrC := make(chan *RunningResult, len(peers))
	worker := metadata.GetPeerWorker(s.runtimeType)

	for _, pd := range peers {
		d := worker.PeerCreateData(pd)
		go func(data runtime.ServiceMetadata) {
			rr := &RunningResult{DataID: data.DataID()}
			err := s.runner.Run(ctx, data)
			if err != nil {
				rr.Err = errors.WithMessagef(err, "启动服务失败：DataID=%s", data.DataID())
			} else {
				rr.Msg = "启动成功"
			}
			rrC <- rr
		}(d)
	}
	rrs := make([]*RunningResult, 0, len(peers))
	for rr := range rrC {
		rrs = append(rrs, rr)
		if len(rrs) == len(peers) {
			break
		}
	}
	msg, err := countResult(rrs)
	if err != nil {
		return err
	}
	fmt.Println(msg)
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
