package agentExec

import (
	"github.com/aaronchen2k/deeptest/internal/agent/exec/domain"
	"github.com/aaronchen2k/deeptest/internal/agent/exec/utils/exec"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
)

type ProcessorLogic struct {
	ID uint `json:"id" yaml:"id"`
	ProcessorEntityBase

	Expression string `json:"expression" yaml:"expression"`
}

func (entity ProcessorLogic) Run(processor *Processor, session *Session) (log domain.Result, err error) {
	processor.Result = domain.Result{
		ID:                entity.ProcessorID,
		Name:              entity.Name,
		ProcessorCategory: entity.ProcessorCategory,
		ProcessorType:     entity.ProcessorType,
		ParentId:          entity.ParentID,
	}

	typ := entity.ProcessorType
	pass := false

	if typ == consts.ProcessorLogicIf {
		var result interface{}
		result, err = EvaluateGovaluateExpressionByScope(entity.Expression, entity.ProcessorID)
		if err != nil {
			pass = false
		} else {
			pass = result.(bool)
		}

	} else if typ == consts.ProcessorLogicElse {
		brother, ok := getPreviousBrother(*processor)
		if ok && brother.Result.ResultStatus != consts.Pass {
			pass = true
		}
	}

	processor.Result.ResultStatus, processor.Result.Summary = getResultStatus(pass)
	processor.Parent.Result.Children = append(processor.Parent.Result.Children, &processor.Result)

	exec.SendExecMsg(processor.Result, session.WsMsg)

	return
}
