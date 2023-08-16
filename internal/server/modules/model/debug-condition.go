package model

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
)

type DebugPreCondition struct {
	BaseModel

	DebugInterfaceId    uint `gorm:"default:0" json:"debugInterfaceId"`
	EndpointInterfaceId uint `gorm:"default:0" json:"endpointInterfaceId"`

	EntityType consts.ConditionType `json:"entityType"`
	EntityId   uint                 `json:"entityId"`
	UsedBy     consts.UsedBy        `json:"usedBy"`

	Name string `json:"name"`
	Desc string `json:"desc"`
	Ordr int    `json:"ordr"`
}

func (DebugPreCondition) TableName() string {
	return "biz_debug_condition_pre"
}

type DebugPostCondition struct {
	BaseModel

	DebugInterfaceId    uint `gorm:"default:0" json:"debugInterfaceId"`
	EndpointInterfaceId uint `gorm:"default:0" json:"endpointInterfaceId"`

	EntityType consts.ConditionType `json:"entityType"`
	EntityId   uint                 `json:"entityId"`
	UsedBy     consts.UsedBy        `json:"usedBy"`

	Name string `json:"name"`
	Desc string `json:"desc"`
	Ordr int    `json:"ordr"`
}

func (DebugPostCondition) TableName() string {
	return "biz_debug_condition_post"
}

type DebugConditionExtractor struct {
	BaseModel

	domain.ExtractorBase
}

func (DebugConditionExtractor) TableName() string {
	return "biz_debug_condition_extractor"
}

type DebugConditionCheckpoint struct {
	BaseModel

	domain.CheckpointBase
}

func (DebugConditionCheckpoint) TableName() string {
	return "biz_debug_condition_checkpoint"
}

type DebugConditionScript struct {
	BaseModel

	domain.ScriptBase
}

func (DebugConditionScript) TableName() string {
	return "biz_debug_condition_script"
}

type DebugConditionCookie struct {
	BaseModel

	domain.CookieBase
}

func (DebugConditionCookie) TableName() string {
	return "biz_debug_condition_cookie"
}

type DebugConditionResponseDefine struct {
	BaseModel
	domain.ResponseDefineBase
	Disabled bool `json:"disabled,omitempty" gorm:"default:false"`
}

func (DebugConditionResponseDefine) TableName() string {
	return "biz_debug_condition_response_define"
}
