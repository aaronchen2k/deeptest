package service

import (
	"github.com/aaronchen2k/deeptest/internal/agent/exec"
	"github.com/aaronchen2k/deeptest/internal/agent/exec/domain"
	"github.com/aaronchen2k/deeptest/internal/agent/exec/utils/exec"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/kataras/iris/v12/websocket"
)

type PlanService struct {
	RemoteService   *RemoteService   `inject:""`
	ScenarioService *ScenarioService `inject:""`
}

func (s *PlanService) ExecPlan(req *agentExec.PlanExecReq, wsMsg *websocket.Message) (err error) {
	consts.ServerUrl = req.ServerUrl
	consts.ServerToken = req.Token

	planExecObj := s.RemoteService.GetPlanToExec(req)

	// start msg
	execUtils.SendStartMsg(wsMsg)

	// execution
	var results = agentDomain.PlanExecResult{
		ID: req.PlanId,
	}
	for _, scenario := range planExecObj.Scenarios {
		session, _ := s.ScenarioService.Exec(&scenario, wsMsg)
		results.Scenarios[scenario.RootProcessor.ScenarioId] = session.RootProcessor.Result
	}

	// submit result
	s.RemoteService.SubmitPlanResult(results, req.PlanId, req.ServerUrl, req.Token)
	s.sendSubmitResult(req.PlanId, wsMsg)

	// end msg
	execUtils.SendEndMsg(wsMsg)

	return
}

func (s *PlanService) CancelAndSendMsg(scenarioId int, wsMsg websocket.Message) (err error) {
	execUtils.SendCancelMsg(wsMsg)
	return
}

func (s *PlanService) RestoreEntityFromRawAndSetParent(root *agentExec.Processor) (err error) {
	processors := make([]*agentExec.Processor, 0)

	agentExec.GetProcessorList(root, &processors)

	processorMap := map[uint]*agentExec.Processor{}
	for _, processor := range processors {
		processorMap[processor.ID] = processor

		processor.RestoreEntity()
	}

	for _, obj := range processorMap {
		obj.Parent = processorMap[obj.ParentId]
	}

	return
}

func (s *PlanService) sendSubmitResult(rootId uint, wsMsg *websocket.Message) (err error) {
	result := agentDomain.ScenarioExecResult{
		ID:       -3,
		ParentId: int(rootId),
		Name:     "提交执行结果成功",
		//Summary:  fmt.Sprintf("错误：%s", err.Error()),
	}
	execUtils.SendExecMsg(result, wsMsg)

	return
}
