package agentExec

import (
	"encoding/json"
	"fmt"
	agentDomain "github.com/aaronchen2k/deeptest/internal/agent/exec/domain"
	execUtils "github.com/aaronchen2k/deeptest/internal/agent/exec/utils/exec"
	ptlog "github.com/aaronchen2k/deeptest/internal/agent/performance/pkg/log"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	commonUtils "github.com/aaronchen2k/deeptest/pkg/lib/comm"
	"github.com/aaronchen2k/deeptest/pkg/lib/log"
	"github.com/kataras/iris/v12"
	"runtime/debug"
	"time"
)

type Processor struct {
	ProcessorBase
	Entity  IProcessorEntity `json:"entity"`
	Disable bool             `json:"disable"`
}

type ProcessorMsg struct {
	ProcessorBase
}

type ProcessorBase struct {
	ID uint `json:"id"`

	Name     string            `json:"name"`
	Comments string            `json:"comments"`
	Method   consts.HttpMethod `json:"method" yaml:"method"`

	ParentId   uint `json:"parentId"`
	ScenarioId uint `json:"scenarioId"`
	ProjectId  uint `json:"projectId"`
	UseID      uint `json:"useId"`

	EntityCategory      consts.ProcessorCategory `json:"entityCategory"`
	EntityType          consts.ProcessorType     `json:"entityType"`
	EntityId            uint                     `json:"entityId"`
	EndpointInterfaceId uint                     `json:"endpointInterfaceId"`

	Ordr      int             `json:"ordr"`
	Children  []*Processor    `json:"children"`
	Slots     iris.Map        `json:"slots"`
	IsDir     bool            `json:"isDir"`
	EntityRaw json.RawMessage `json:"entityRaw"`

	Parent                *Processor                      `json:"-"`
	Result                *agentDomain.ScenarioExecResult `json:"result"`
	ProcessorInterfaceSrc consts.ProcessorInterfaceSrc    `json:"processorInterfaceSrc"`

	Round string `json:"round"`

	Session ExecSession `json:"-"`
}

func (p *Processor) Run(session *ExecSession) (err error) {
	_logUtils.Infof("%d - %s %s", p.ID, p.Name, p.EntityType)

	select {
	case <-session.Ctx.Done():
		break
	default:
	}

	session.CurrScenarioProcessorId = p.ID
	session.CurrScenarioProcessor = p

	//每个执行器延迟0.1秒，防止发送ws消息过快，导致前端消息错误
	time.Sleep(100 * time.Microsecond)
	if !p.Disable && p.Entity != nil {
		p.Entity.Run(p, session)
	}

	return
}

func (p *Processor) Error(s *ExecSession, err interface{}) {
	detail := map[string]interface{}{}

	if p.Result.Detail != "" {
		commonUtils.JsonDecode(p.Result.Detail, &detail)
	}

	detail["exception"] = fmt.Sprintf("错误：%v", err)
	p.Result.Detail = commonUtils.JsonEncode(detail)

	msg := fmt.Sprintf("err=%v\n stack=%s\n", err, string(debug.Stack()))
	_logUtils.Errorf(msg)
	ptlog.Logf(msg)

	p.AddResultToParent()
	execUtils.SendExecMsg(p.Result, consts.Processor, s.WsMsg)

	//execUtils.SendExceptionMsg(s.WsMsg)
}

func (p *Processor) AddResultToParent() (err error) {
	p.Parent.Result.Children = append(p.Parent.Result.Children, p.Result)
	return
}

func (p *Processor) RestoreEntity() (err error) {
	bytes, err := p.EntityRaw.MarshalJSON()

	switch p.EntityCategory {
	case consts.ProcessorInterface:
		ret := ProcessorInterface{}
		json.Unmarshal(bytes, &ret)
		p.Entity = ret

	case consts.ProcessorRoot:
		ret := ProcessorRoot{}
		json.Unmarshal(bytes, &ret)
		p.Entity = ret

	case consts.ProcessorGroup:
		ret := ProcessorGroup{}
		json.Unmarshal(bytes, &ret)
		p.Entity = ret

	case consts.ProcessorLogic:
		ret := ProcessorLogic{}
		json.Unmarshal(bytes, &ret)
		p.Entity = ret

	case consts.ProcessorLoop:
		ret := ProcessorLoop{}
		json.Unmarshal(bytes, &ret)
		p.Entity = ret

	case consts.ProcessorVariable:
		ret := ProcessorVariable{}
		json.Unmarshal(bytes, &ret)
		p.Entity = ret

	case consts.ProcessorTimer:
		ret := ProcessorTimer{}
		json.Unmarshal(bytes, &ret)
		p.Entity = ret

	case consts.ProcessorPrint:
		ret := ProcessorPrint{}
		json.Unmarshal(bytes, &ret)
		p.Entity = ret

	case consts.ProcessorCookie:
		ret := ProcessorCookie{}
		json.Unmarshal(bytes, &ret)
		p.Entity = ret

	case consts.ProcessorAssertion:
		ret := ProcessorAssertion{}
		json.Unmarshal(bytes, &ret)
		p.Entity = ret

	case consts.ProcessorData:
		ret := ProcessorData{}
		json.Unmarshal(bytes, &ret)
		p.Entity = ret

	case consts.ProcessorCustomCode:
		ret := ProcessorCustomCode{}
		json.Unmarshal(bytes, &ret)
		p.Entity = ret

	//case consts.ProcessorPerformanceGoal:
	//	ret := ProcessorPerformanceGoal{}
	//	json.Unmarshal(bytes, &ret)
	//	p.Entity = ret
	//
	//case consts.ProcessorPerformanceRunner:
	//	ret := ProcessorPerformanceRunner{}
	//	json.Unmarshal(bytes, &ret)
	//	p.Entity = ret
	//
	//case consts.ProcessorPerformanceScenario:
	//	ret := ProcessorPerformanceScenario{}
	//	json.Unmarshal(bytes, &ret)
	//	p.Entity = ret
	//
	//case consts.ProcessorPerformanceRendezvous:
	//	ret := ProcessorPerformanceRendezvous{}
	//	json.Unmarshal(bytes, &ret)
	//	p.Entity = ret

	default:
	}

	return
}

func (p *Processor) AppendNewChildProcessor(category consts.ProcessorCategory, typ consts.ProcessorType) (child Processor) {
	child = Processor{
		ProcessorBase: ProcessorBase{
			EntityCategory: category,
			EntityType:     typ,
			Parent:         p,
			ParentId:       p.ID,
		},
	}

	child.Result = &agentDomain.ScenarioExecResult{
		ProcessorCategory: child.EntityCategory,
		ProcessorType:     child.EntityType,
		ParentId:          int(p.ID),
	}

	return
}
