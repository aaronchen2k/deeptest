package runnerExec

import (
	"context"
	performanceUtils "github.com/aaronchen2k/deeptest/internal/agent/exec/utils/performance"
	ptlog "github.com/aaronchen2k/deeptest/internal/agent/performance/pkg/log"
	ptProto "github.com/aaronchen2k/deeptest/internal/agent/performance/proto"
	_logUtils "github.com/aaronchen2k/deeptest/pkg/lib/log"
	_stringUtils "github.com/aaronchen2k/deeptest/pkg/lib/string"
	"sync"
	"time"
)

type ConstantVuGenerator struct {
}

func (g ConstantVuGenerator) Run(execCtx context.Context) (err error) {
	execParams := performanceUtils.GetExecParamsInCtx(execCtx)
	ptlog.Logf("Constant Generator run, execParams: %s", _stringUtils.ToJsonStr(execParams))

	var wgVus sync.WaitGroup

	target := performanceUtils.GetVuNumbByWeight(execParams.Target, execParams.Weight)

	for i := 1; i <= target; i++ {
		vuCtx, _ := context.WithCancel(execCtx)
		if execParams.GoalDuration > 0 { // control exec time
			vuCtx, _ = context.WithTimeout(execCtx, time.Duration(execParams.GoalDuration)*time.Second)
		}

		wgVus.Add(1)

		result := ptProto.PerformanceExecResp{
			Timestamp:  time.Now().UnixMilli(),
			RunnerId:   execParams.RunnerId,
			RunnerName: execParams.RunnerName,
			Room:       execParams.Room,

			VuCount: 1,
		}
		execParams.Sender.Send(result)

		index := i
		go func() {
			defer wgVus.Done()

			//execParams.VuNo = index
			ExecScenarioWithVu(vuCtx, index)

			ptlog.Logf("vu %d completed", index)
		}()

		select {
		case <-vuCtx.Done():
			_logUtils.Debug("<<<<<<< stop stages")
			goto Label_END_STAGES

		default:
		}
	}

	// wait
	wgVus.Wait()

	ptlog.Logf("all vus completed")

Label_END_STAGES:

	return
}
