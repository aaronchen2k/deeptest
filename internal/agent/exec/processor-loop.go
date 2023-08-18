package agentExec

import (
	"fmt"
	"github.com/aaronchen2k/deeptest/internal/agent/exec/domain"
	agentUtils "github.com/aaronchen2k/deeptest/internal/agent/exec/utils"
	"github.com/aaronchen2k/deeptest/internal/agent/exec/utils/exec"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	commonUtils "github.com/aaronchen2k/deeptest/pkg/lib/comm"
	logUtils "github.com/aaronchen2k/deeptest/pkg/lib/log"
	uuid "github.com/satori/go.uuid"
	"time"
)

type ProcessorLoop struct {
	ID uint `json:"id" yaml:"id"`
	ProcessorEntityBase

	Times        int    `json:"times" yaml:"times"` // time
	Range        string `json:"range" yaml:"range"` // range
	List         string `json:"list" yaml:"list"`   // in
	Step         string `json:"step" yaml:"step"`
	IsRand       bool   `json:"isRand" yaml:"isRand"`
	VariableName string `json:"variableName" yaml:"variableName"`

	UntilExpression   string `json:"untilExpression" yaml:"untilExpression"` // until
	BreakIfExpression string `json:"breakIfExpression" yaml:"breakIfExpression"`
}

func (entity ProcessorLoop) Run(processor *Processor, session *Session) (err error) {
	logUtils.Infof("loop entity")

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

	processor.Result.Detail = commonUtils.JsonEncode(entity)
	execUtils.SendExecMsg(*processor.Result, session.WsMsg)

	if entity.ProcessorType == consts.ProcessorLoopBreak {
		processor.Result.WillBreak, processor.Result.Summary, processor.Result.Detail = entity.getBeak()

		processor.AddResultToParent()
		if processor.Result.WillBreak {
			execUtils.SendExecMsg(*processor.Result, session.WsMsg)
		}

		return
	}
	processor.Result.Iterator, processor.Result.Summary = entity.getIterator()

	if entity.ProcessorType == consts.ProcessorLoopUntil {
		entity.runLoopUntil(session, processor, processor.Result.Iterator)
	} else {
		entity.runLoopItems(session, processor, processor.Result.Iterator)
	}

	processor.AddResultToParent()

	endTime := time.Now()
	processor.Result.EndTime = &endTime

	return
}

func (entity *ProcessorLoop) runLoopItems(session *Session, processor *Processor, iterator agentDomain.ExecIterator) (err error) {
	for index, item := range iterator.Items {
		if ForceStopExec {
			break
		}
		if DemoTestSite != "" && index > 2 {
			break
		}
		/*
			msg := agentDomain.ScenarioExecResult{
				ParentId:          int(processor.ID),
				Summary:           fmt.Sprintf("%d. %s为%v", index+1, iterator.VariableName, item),
				Name:              "循环变量",
				ProcessorCategory: consts.ProcessorPrint,
			}
			execUtils.SendExecMsg(msg, session.WsMsg)
		*/

		SetVariable(entity.ProcessorID, iterator.VariableName, item, consts.Public)
		round := ""
		for _, child := range processor.Children {

			if ForceStopExec {
				break
			}
			if child.Disable {
				continue
			}
			//执行轮次
			if round == "" {
				if entity.ProcessorType == consts.ProcessorLoopTime {
					round = fmt.Sprintf("第 %v 轮", index+1)
				} else {
					round = fmt.Sprintf("第 %v 轮，%v = %v", index+1, iterator.VariableName, item)
				}

				child.Round = round
			}

			(*child).Run(session)

			if child.Result.WillBreak {
				logUtils.Infof("break")
				goto LABEL
			}
		}
	}
LABEL:

	return
}

func (entity *ProcessorLoop) runLoopUntil(session *Session, processor *Processor, iterator agentDomain.ExecIterator) (err error) {
	expression := iterator.UntilExpression

	index := 0
	for {
		if ForceStopExec {
			break
		}
		if DemoTestSite != "" && index > 2 {
			break
		}
		index += 1

		/*
			msg := agentDomain.ScenarioExecResult{
				ParentId: int(processor.ID),
				Summary:  fmt.Sprintf("%d. ", index),
			}
			execUtils.SendExecMsg(msg, session.WsMsg)
		*/

		result, err := EvaluateGovaluateExpressionByProcessorScope(expression, entity.ProcessorID)
		pass, ok := result.(bool)
		if err != nil || !ok || pass {
			childBreakProcessor := processor.AppendNewChildProcessor(consts.ProcessorLoop, consts.ProcessorLoopBreak)
			childBreakProcessor.Result.WillBreak = true
			childBreakProcessor.Result.Summary = fmt.Sprintf("条件%s满足，跳出循环。", expression)

			childBreakProcessor.AddResultToParent()
			execUtils.SendExecMsg(*childBreakProcessor.Result, session.WsMsg)

			break
		}

		round := ""

		for _, child := range processor.Children {
			if ForceStopExec {
				break
			}
			if child.Disable {
				continue
			}

			if round == "" {
				round = fmt.Sprintf("第 %v 轮", index+1)
				child.Round = round
			}

			(*child).Run(session)

			if child.Result.WillBreak {
				logUtils.Infof("break")
				goto LABEL
			}
		}

		if index >= consts.MaxLoopTimeForInterfaceTest {
			logUtils.Infof("break for reach MaxLoopTimeForInterfaceTest")
			goto LABEL
		}
	}
LABEL:

	return
}

func (entity *ProcessorLoop) getBeak() (ret bool, msg string, detailStr string) {
	breakFrom := entity.ParentID
	breakIfExpress := entity.BreakIfExpression

	result, err := EvaluateGovaluateExpressionByProcessorScope(breakIfExpress, entity.ProcessorID)
	ret, ok := result.(bool)
	pass := false
	if err == nil && ok && ret {
		breakMap.Store(breakFrom, true)
		msg = "真"
		pass = true
	} else {
		msg = "假"
	}

	//detail := map[string]interface{}{"表达式": breakIfExpress + " 为 " + msg}
	detail := map[string]interface{}{"expression": breakIfExpress, "result": pass}
	detailStr = commonUtils.JsonEncode(detail)

	return
}

func (entity *ProcessorLoop) getIterator() (iterator agentDomain.ExecIterator, msg string) {
	if entity.ID == 0 {
		msg = "执行前请先配置处理器。"
		return
	}

	if entity.ProcessorType == consts.ProcessorLoopTime {
		iterator, _ = entity.GenerateLoopTimes()
		msg = fmt.Sprintf("迭代\"%d\"次。", entity.Times)
	} else if entity.ProcessorType == consts.ProcessorLoopIn {
		iterator, _ = entity.GenerateLoopList()
		msg = fmt.Sprintf("\"%s\"。", entity.List)
	} else if entity.ProcessorType == consts.ProcessorLoopRange {
		iterator, _ = entity.GenerateLoopRange()
		msg = fmt.Sprintf("\"%s\"。", entity.Range)
	} else if entity.ProcessorType == consts.ProcessorLoopUntil {
		iterator.UntilExpression = entity.UntilExpression
		msg = fmt.Sprintf("\"%s\"。", entity.UntilExpression)
	}

	iterator.VariableName = entity.VariableName

	return
}

func (entity *ProcessorLoop) GenerateLoopTimes() (ret agentDomain.ExecIterator, err error) {
	if entity.Times > 0 {
		for i := 0; i < entity.Times; i++ {
			ret.Items = append(ret.Items, i+1)
		}
	}

	ret.DataType = consts.Int

	return
}
func (entity *ProcessorLoop) GenerateLoopRange() (ret agentDomain.ExecIterator, err error) {
	start, end, step, precision, typ, err := agentUtils.GetRange(entity.Range, entity.Step)
	if err == nil {
		ret.DataType = typ
		ret.Items, _ = agentUtils.GenerateRangeItems(start, end, step, precision, entity.IsRand, typ)
	}

	return
}
func (entity *ProcessorLoop) GenerateLoopList() (ret agentDomain.ExecIterator, err error) {
	ret.Items, ret.DataType, err = agentUtils.GenerateListItems(entity.List, entity.IsRand)

	return
}
