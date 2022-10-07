package domain

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"time"
)

type ExecVariable struct {
	Id         uint        `json:"id"`
	Name       string      `json:"name"`
	Value      interface{} `json:"value"`
	Expression string      `json:"expression"`

	InterfaceId uint                  `json:"interfaceId"`
	Scope       consts.ExtractorScope `json:"scope"`
}

type ExecCookie struct {
	Id    uint        `json:"id"`
	Name  string      `json:"name"`
	Value interface{} `json:"value"`

	Domain     string     `json:"domain"`
	ExpireTime *time.Time `json:"expireTime"`
}

type Extractor struct {
	Src  consts.ExtractorSrc  `json:"src"`
	Type consts.ExtractorType `json:"type"`
	Key  string               `json:"key"`

	Expression string `json:"expression"`
	Prop       string `json:"prop"`

	BoundaryStart    string `json:"boundaryStart"`
	BoundaryEnd      string `json:"boundaryEnd"`
	BoundaryIndex    int    `json:"boundaryIndex"`
	BoundaryIncluded bool   `json:"boundaryIncluded"`

	Variable string                `json:"variable"`
	Scope    consts.ExtractorScope `json:"scope"`

	Result      string `json:"result"`
	InterfaceId uint   `json:"interfaceId"`

	Disabled bool `json:"disabled"`
}
type Checkpoint struct {
	Type consts.CheckpointType `json:"type"`

	Expression        string `json:"expression"`
	ExtractorVariable string `json:"extractorVariable"`

	Operator     consts.ComparisonOperator `json:"operator"`
	Value        string                    `json:"value"`
	ActualResult string                    `json:"actualResult"`

	ResultStatus consts.ResultStatus `json:"resultStatus"`
	InterfaceId  uint                `json:"interfaceId"`

	Disabled bool `json:"disabled"`
}

type ExecIterator struct {
	ProcessorCategory consts.ProcessorCategory
	ProcessorType     consts.ProcessorType
	VariableName      string `json:"variableName,omitempty"`

	// loop range
	Items    []interface{}            `json:"items"`
	Data     []map[string]interface{} `json:"data"`
	DataType consts.DataType          `json:"dataType"`

	// loop condition
	UntilExpression string `json:"untilExpression"`
}

type ExecOutput struct {
	// logic if, else
	Pass bool `json:"pass,omitempty"`

	// loop - times
	Times int `json:"times,omitempty"`
	// loop util
	Expression string `json:"times,omitempty"`
	// loop in
	List string `json:"list,omitempty"`
	// loop - range
	Range      string          `json:"range,omitempty"`
	RangeStart interface{}     `json:"rangeStart,omitempty"`
	RangeEnd   interface{}     `json:"rangeEnd,omitempty"`
	RangeType  consts.DataType `json:"rangeType,omitempty"`
	// loop break
	BreakFrom uint `json:"breakFrom,omitempty"`

	// timer
	SleepTime int `json:"sleepTime"`

	// data
	Url          string `json:"url"`
	RepeatTimes  int    `json:"repeatTimes,omitempty"`
	IsLoop       int    `json:"isLoop,omitempty"`
	IsRand       bool   `json:"isRand,omitempty"`
	IsOnce       bool   `json:"isOnce,omitempty"`
	VariableName string `json:"variableName,omitempty"`

	// extractor
	Src  consts.ExtractorSrc  `json:"src"`
	Type consts.ExtractorType `json:"type"`
	//Expression string `json:"expression"`
	BoundaryStart    string `json:"boundaryStart"`
	BoundaryEnd      string `json:"boundaryEnd"`
	BoundaryIndex    int    `json:"boundaryIndex"`
	BoundaryIncluded bool   `json:"boundaryIncluded"`
	Variable         string `json:"variable"`

	// variable
	VariableValue interface{} `json:"variableValue"`

	// common
	Msg string `json:"msg,omitempty"`
}
