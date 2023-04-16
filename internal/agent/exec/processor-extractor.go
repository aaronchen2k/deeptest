package agentExec

import (
	"encoding/json"
	"fmt"
	v1 "github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/agent/exec/domain"
	"github.com/aaronchen2k/deeptest/internal/agent/exec/utils/exec"
	"github.com/aaronchen2k/deeptest/internal/agent/exec/utils/query"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	httpHelper "github.com/aaronchen2k/deeptest/internal/pkg/helper/http"
	logUtils "github.com/aaronchen2k/deeptest/pkg/lib/log"
	"strings"
	"time"
)

type ProcessorExtractor struct {
	ID uint `json:"id" yaml:"id"`
	ProcessorEntityBase

	Src  consts.ExtractorSrc  `json:"src"`
	Type consts.ExtractorType `json:"type"`
	Key  string               `json:"key"` // form header

	Expression string `json:"expression"`
	//NodeProp       string `json:"prop"`

	BoundaryStart    string `json:"boundaryStart"`
	BoundaryEnd      string `json:"boundaryEnd"`
	BoundaryIndex    int    `json:"boundaryIndex"`
	BoundaryIncluded bool   `json:"boundaryIncluded"`

	Variable string `json:"variable"`

	Result      string `json:"result"`
	InterfaceID uint   `json:"interfaceID"`

	Disabled bool `json:"disabled"`
}

func (entity ProcessorExtractor) Run(processor *Processor, session *Session) (err error) {
	logUtils.Infof("extractor entity")

	startTime := time.Now()
	processor.Result = &agentDomain.ScenarioExecResult{
		ID:                int(entity.ProcessorID),
		Name:              entity.Name,
		ProcessorCategory: entity.ProcessorCategory,
		ProcessorType:     entity.ProcessorType,
		StartTime:         &startTime,
		ParentId:          int(entity.ParentID),
	}

	brother, ok := getPreviousBrother(*processor)
	if !ok || brother.EntityType != consts.ProcessorInterfaceDefault {
		processor.Result.Summary = fmt.Sprintf("先前节点不是接口，无法应用提取器。")
		processor.AddResultToParent()
		execUtils.SendExecMsg(*processor.Result, session.WsMsg)
		return
	}

	resp := v1.DebugResponse{}
	json.Unmarshal([]byte(brother.Result.RespContent), &resp)

	entity.Src = consts.Body
	entity.Type = getExtractorTypeForProcessor(entity.ProcessorType)

	err = ExtractValue(&entity, resp)
	if err != nil {
		processor.Result.Summary = fmt.Sprintf("%s提取器解析错误 %s。", entity.ProcessorType, err.Error())
		processor.AddResultToParent()
		execUtils.SendExecMsg(*processor.Result, session.WsMsg)
		return
	}

	SetVariable(processor.ParentId, entity.Variable, entity.Result, consts.Public) // set in parent scope

	processor.Result.Summary = fmt.Sprintf("将结果\"%v\"赋予变量\"%s\"。", entity.Result, entity.Variable)
	processor.AddResultToParent()
	execUtils.SendExecMsg(*processor.Result, session.WsMsg)

	endTime := time.Now()
	processor.Result.EndTime = &endTime

	return
}

func ExtractValue(extractor *ProcessorExtractor, resp v1.DebugResponse) (err error) {
	if extractor.Disabled {
		extractor.Result = ""
		return
	}

	if extractor.Src == consts.Header {
		for _, h := range resp.Headers {
			if h.Name == extractor.Key {
				extractor.Result = h.Value
				break
			}
		}
	} else {
		if httpHelper.IsJsonContent(resp.ContentType.String()) && extractor.Type == consts.JsonQuery {
			extractor.Result = queryUtils.JsonQuery(resp.Content, extractor.Expression)

		} else if httpHelper.IsHtmlContent(resp.ContentType.String()) && extractor.Type == consts.HtmlQuery {
			extractor.Result = queryUtils.HtmlQuery(resp.Content, extractor.Expression)

		} else if httpHelper.IsXmlContent(resp.ContentType.String()) && extractor.Type == consts.XmlQuery {
			extractor.Result = queryUtils.XmlQuery(resp.Content, extractor.Expression)

		} else if extractor.Type == consts.Boundary {
			extractor.Result = queryUtils.BoundaryQuery(resp.Content, extractor.BoundaryStart, extractor.BoundaryEnd,
				extractor.BoundaryIndex, extractor.BoundaryIncluded)
		}
	}

	extractor.Result = strings.TrimSpace(extractor.Result)

	return
}

func getExtractorTypeForProcessor(processorType consts.ProcessorType) (ret consts.ExtractorType) {
	if processorType == consts.ProcessorExtractorBoundary {
		ret = consts.Boundary
	} else if processorType == consts.ProcessorExtractorJsonQuery {
		ret = consts.JsonQuery
	} else if processorType == consts.ProcessorExtractorHtmlQuery {
		ret = consts.HtmlQuery
	} else if processorType == consts.ProcessorExtractorXmlQuery {
		ret = consts.XmlQuery
	}

	return
}
