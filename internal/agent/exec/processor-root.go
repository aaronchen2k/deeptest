package agentExec

import (
	"github.com/aaronchen2k/deeptest/internal/agent/exec/domain"
	"github.com/aaronchen2k/deeptest/internal/agent/exec/utils/exec"
	logUtils "github.com/aaronchen2k/deeptest/pkg/lib/log"
	"time"
)

type ProcessorRoot struct {
	ID uint `json:"id" yaml:"id"`
	ProcessorEntityBase
}

func (entity ProcessorRoot) Run(processor *Processor, session *Session) (err error) {
	logUtils.Infof("root entity")

	startTime := time.Now()
	processor.Result = &agentDomain.ScenarioExecResult{
		ID:                int(entity.ProcessorID),
		Name:              session.Name,
		ProcessorCategory: entity.ProcessorCategory,
		ProcessorType:     entity.ProcessorType,
		StartTime:         &startTime,
		ParentId:          int(entity.ParentID),
		ScenarioId:        processor.ScenarioId,
	}

	execUtils.SendExecMsg(*processor.Result, session.WsMsg)

	for _, child := range processor.Children {
		child.Run(session)
	}

	endTime := time.Now()
	processor.Result.EndTime = &endTime

	return
}
