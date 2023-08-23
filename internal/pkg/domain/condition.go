package domain

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
)

type ExtractorBase struct {
	Src  consts.ExtractorSrc  `json:"src"`
	Type consts.ExtractorType `json:"type"`
	Key  string               `json:"key"`

	Expression string `gorm:"default:''" json:"expression"`
	Prop       string `json:"prop"`

	BoundaryStart    string `gorm:"default:''" json:"boundaryStart"`
	BoundaryEnd      string `gorm:"default:''" json:"boundaryEnd"`
	BoundaryIndex    int    `json:"boundaryIndex"`
	BoundaryIncluded bool   `json:"boundaryIncluded"`

	Variable string                `gorm:"default:''" json:"variable"`
	Scope    consts.ExtractorScope `json:"scope" gorm:"default:public"`

	Default string `gorm:"default:''" json:"default"` // for cookie

	Result       string              `json:"result"`
	ResultStatus consts.ResultStatus `json:"resultStatus"`
	ResultMsg    string              `json:"resultMsg"`

	ConditionId         uint                 `json:"conditionId"`
	ConditionEntityId   uint                 `gorm:"-" json:"conditionEntityId"`   // refer to po id in domain object
	ConditionEntityType consts.ConditionType `gorm:"-" json:"conditionEntityType"` // for log only
	InvokeId            uint                 `json:"invokeId"`                     // for log only

	Disabled bool `json:"disabled"`
}

func (condition ExtractorBase) GetType() consts.ConditionType {
	return consts.ConditionTypeExtractor
}

type CheckpointBase struct {
	Type consts.CheckpointType `json:"type"`

	Expression        string `json:"expression"`
	ExtractorVariable string `json:"extractorVariable"`

	Operator     consts.ComparisonOperator `json:"operator"`
	Value        string                    `json:"value"`
	ActualResult string                    `json:"actualResult"`

	ResultStatus consts.ResultStatus `json:"resultStatus"`
	ResultMsg    string              `json:"resultMsg"`

	ConditionId         uint                 `json:"conditionId"`
	ConditionEntityId   uint                 `gorm:"-" json:"conditionEntityId"`   // refer to entity po id in domain object
	ConditionEntityType consts.ConditionType `gorm:"-" json:"conditionEntityType"` // for log only
	InvokeId            uint                 `json:"invokeId"`                     // for log only

	Disabled bool `json:"disabled"`
}

func (condition CheckpointBase) GetType() consts.ConditionType {
	return consts.ConditionTypeCheckpoint
}

type ScriptBase struct {
	ConditionSrc consts.ConditionSrc `json:"conditionSrc"`

	Content string `gorm:"type:longtext;" json:"content"`

	Output       string              `gorm:"type:longtext;" json:"output"`
	ResultStatus consts.ResultStatus `json:"resultStatus"`
	ResultMsg    string              `json:"resultMsg"`

	ConditionId         uint                 `json:"conditionId"`
	ConditionEntityId   uint                 `gorm:"-" json:"conditionEntityId"`   // refer to po id in domain object
	ConditionEntityType consts.ConditionType `gorm:"-" json:"conditionEntityType"` // for log only
	InvokeId            uint                 `json:"invokeId"`                     // for log only

	Disabled bool `json:"disabled"`

	VariableSettings []ExecVariable `gorm:"-" json:"variableSettings"`
}

func (condition ScriptBase) GetType() consts.ConditionType {
	return consts.ConditionTypeScript
}

type ResponseDefineBase struct {
	ResponseCode string   `json:"responseCode"`
	Schema       string   `gorm:"-" json:"schema"`
	Codes        []string `gorm:"-" json:"codes"`
	Code         string   `json:"code"`

	Output       string              `gorm:"type:longtext;" json:"output"`
	ResultStatus consts.ResultStatus `json:"resultStatus"`
	ResultMsg    string              `json:"resultMsg"`

	ConditionId         uint                 `json:"conditionId"`
	ConditionEntityId   uint                 `gorm:"-" json:"conditionEntityId"`   // refer to po id in domain object
	ConditionEntityType consts.ConditionType `gorm:"-" json:"conditionEntityType"` // for log only
	InvokeId            uint                 `json:"invokeId"`                     // for log only
	MediaType           string               `json:"mediaType"`
	Disabled            bool                 `json:"disabled"`
	Component           string               `gorm:"-" json:"component"`
}

func (condition ResponseDefineBase) GetType() consts.ConditionType {
	return consts.ConditionTypeResponseDefine
}
