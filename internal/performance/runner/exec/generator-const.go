package runnerExec

import (
	"context"
	ptlog "github.com/aaronchen2k/deeptest/internal/performance/pkg/log"
	ptProto "github.com/aaronchen2k/deeptest/internal/performance/proto"
	"github.com/aaronchen2k/deeptest/internal/performance/runner/metrics"
	_logUtils "github.com/aaronchen2k/deeptest/pkg/lib/log"
	"sync"
	"time"
)

type ConstantVuGenerator struct {
}

func (g ConstantVuGenerator) Run(execCtx context.Context, sender metrics.MessageSender) (err error) {
	execParams := getExecParamsInCtx(execCtx)

	var wgVus sync.WaitGroup

	target := getVuNumbByWeight(execParams.Target, execParams.Weight)

	for i := 1; i <= target; i++ {
		childCtx := execCtx
		if execParams.Duration > 0 {
			childCtx, _ = context.WithTimeout(execCtx, time.Duration(execParams.Duration)*time.Second)
		}

		wgVus.Add(1)

		result := ptProto.PerformanceExecResp{
			Timestamp:  time.Now().UnixMilli(),
			RunnerId:   execParams.RunnerId,
			RunnerName: execParams.RunnerName,
			Room:       execParams.Room,

			VuCount: 1,
		}
		sender.Send(result)

		index := i
		go func() {
			defer wgVus.Done()
			ExecScenarioWithVu(childCtx, sender, index)

			ptlog.Logf("vu %d completed", index)
		}()

		select {
		case <-childCtx.Done():
			_logUtils.Debug("<<<<<<< stop stages")
			goto Label_END_STAGES

		default:
		}
	}

	// wait all vus completed
	wgVus.Wait()

	ptlog.Log("all vus completed")

Label_END_STAGES:

	return
}