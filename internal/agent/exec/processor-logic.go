package agentExec

import (
	"github.com/aaronchen2k/deeptest/internal/agent/exec/domain"
	"github.com/aaronchen2k/deeptest/internal/agent/exec/utils/exec"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	commonUtils "github.com/aaronchen2k/deeptest/pkg/lib/comm"
	logUtils "github.com/aaronchen2k/deeptest/pkg/lib/log"
	uuid "github.com/satori/go.uuid"
	"time"
)

type ProcessorLogic struct {
	ID uint `json:"id" yaml:"id"`
	ProcessorEntityBase

	Expression string `json:"expression" yaml:"expression"`
}

func (entity ProcessorLogic) Run(processor *Processor, session *ExecSession) (err error) {
	defer func() {
		if errX := recover(); errX != nil {
			processor.Error(session, errX)
		}
	}()
	logUtils.Infof("logic entity")

	startTime := time.Now()
	processor.Result = &agentDomain.ScenarioExecResult{
		ID:                int(entity.ProcessorID),
		Name:              entity.Name,
		ProcessorCategory: entity.ProcessorCategory,
		ProcessorType:     entity.ProcessorType,
		StartTime:         &startTime,
		ParentId:          int(entity.ParentID),
		ScenarioId:        processor.ScenarioId,
		ProcessorId:       processor.ID,
		LogId:             uuid.NewV4(),
		ParentLogId:       processor.Parent.Result.LogId,
		Round:             processor.Round,
	}

	typ := entity.ProcessorType
	pass := false
	detail := map[string]interface{}{"name": entity.Name, "expression": entity.Expression}
	if typ == consts.ProcessorLogicIf {
		result, _ := NewGojaSimple().ExecJsFuncSimple(entity.Expression, session, true)
		pass, _ = result.(bool)

	} else if typ == consts.ProcessorLogicElse {
		brother, ok := getPreviousBrother(*processor)
		if ok && brother.Result.ResultStatus != consts.Pass {
			pass = true
		}
	}

	processor.Result.ResultStatus, processor.Result.Summary = getResultStatus(pass)
	detail["result"] = pass
	processor.Result.Detail = commonUtils.JsonEncode(detail)
	execUtils.SendExecMsg(*processor.Result, consts.Processor, session.ScenarioDebug.WsMsg)
	processor.AddResultToParent()

	executedProcessorIds := map[uint]bool{}
	if pass {
		for _, child := range processor.Children {
			if session.GetForceStop() {
				break
			}
			if child.Disable {
				continue
			}

			executedProcessorIds[child.ID] = true

			child.Run(session)
		}
	}

	endTime := time.Now()
	processor.Result.EndTime = &endTime

	stat := CountSkip(executedProcessorIds, processor.Children, session)
	execUtils.SendStatMsg(stat, session.ScenarioDebug.WsMsg)

	return
}
