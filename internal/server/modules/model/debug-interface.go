package model

import "github.com/aaronchen2k/deeptest/internal/pkg/consts"

type DebugInterface struct {
	BaseModel

	InterfaceBase

	//ReqBodySpec InterfaceReqBodySpec `gorm:"-" json:"basicAuth"`

	Children []*DebugInterface `gorm:"-" json:"children"`

	Params  []DebugInterfaceParam  `gorm:"-" json:"params"`
	Headers []DebugInterfaceHeader `gorm:"-" json:"headers"`
	Cookies []DebugInterfaceCookie `gorm:"-" json:"cookies"`

	BodyFormData       []DebugInterfaceBodyFormDataItem       `gorm:"-" json:"bodyFormData"`
	BodyFormUrlencoded []DebugInterfaceBodyFormUrlEncodedItem `gorm:"-" json:"bodyFormUrlencoded"`

	BasicAuth   DebugInterfaceBasicAuth   `gorm:"-" json:"basicAuth"`
	BearerToken DebugInterfaceBearerToken `gorm:"-" json:"bearerToken"`
	OAuth20     DebugInterfaceOAuth20     `gorm:"-" json:"oauth20"`
	ApiKey      DebugInterfaceApiKey      `gorm:"-" json:"apiKey"`

	InterfaceExtractors  []DebugInterfaceExtractor  `gorm:"-" json:"interfaceExtractors"`
	InterfaceCheckpoints []DebugInterfaceCheckpoint `gorm:"-" json:"interfaceCheckpoints"`
}

func (DebugInterface) TableName() string {
	return "biz_debug_interface"
}

type DebugInterfaceParam struct {
	BaseModel
	InterfaceParamBase
}

func (DebugInterfaceParam) TableName() string {
	return "biz_debug_interface_param"
}

type DebugInterfaceBodyFormDataItem struct {
	BaseModel
	InterfaceBodyFormDataItemBase
}

func (DebugInterfaceBodyFormDataItem) TableName() string {
	return "biz_debug_interface_form_data_item"
}

type DebugInterfaceBodyFormUrlEncodedItem struct {
	BaseModel
	InterfaceBodyFormUrlEncodedItemBase
}

func (DebugInterfaceBodyFormUrlEncodedItem) TableName() string {
	return "biz_debug_interface_form_urlencoded_item"
}

type DebugInterfaceHeader struct {
	BaseModel
	InterfaceHeaderBase
}

func (DebugInterfaceHeader) TableName() string {
	return "biz_debug_interface_header"
}

type DebugInterfaceCookie struct {
	BaseModel
	InterfaceCookieBase
}

func (DebugInterfaceCookie) TableName() string {
	return "biz_debug_interface_cookie"
}

type DebugInterfaceBasicAuth struct {
	BaseModel
	InterfaceBasicAuthBase
}

func (DebugInterfaceBasicAuth) TableName() string {
	return "biz_debug_interface_basic_auth"
}

type DebugInterfaceBearerToken struct {
	BaseModel
	InterfaceBearerTokenBase
}

func (DebugInterfaceBearerToken) TableName() string {
	return "biz_debug_interface_bearer_token"
}

type DebugInterfaceOAuth20 struct {
	BaseModel

	InterfaceOAuth20Base
}

func (DebugInterfaceOAuth20) TableName() string {
	return "biz_debug_interface_oauth20"
}

type DebugInterfaceApiKey struct {
	BaseModel
	InterfaceApiKeyBase
}

func (DebugInterfaceApiKey) TableName() string {
	return "biz_debug_interface_apikey"
}

type DebugInterfaceExtractor struct {
	BaseModel

	UsedBy consts.UsedBy `json:"usedBy"`

	Src  consts.ExtractorSrc  `json:"src"`
	Type consts.ExtractorType `json:"type"`
	Key  string               `json:"key"`

	Expression string `json:"expression"`
	//NodeProp       string `json:"prop"`

	BoundaryStart    string `json:"boundaryStart"`
	BoundaryEnd      string `json:"boundaryEnd"`
	BoundaryIndex    int    `json:"boundaryIndex"`
	BoundaryIncluded bool   `json:"boundaryIncluded"`

	Variable string                `json:"variable"`
	Scope    consts.ExtractorScope `json:"scope" gorm:"default:private"`
	//DisableShare bool                  `json:"disableShare"`

	Result      string `json:"result" gorm:"type:text"`
	InterfaceId uint   `json:"interfaceId"`
	ProcessorId uint   `json:"processorId"`
	ScenarioId  uint   `json:"scenarioId"`

	ProjectId uint `json:"projectId"`
}

func (DebugInterfaceExtractor) TableName() string {
	return "biz_debug_interface_extractor"
}

type DebugInterfaceCheckpoint struct {
	BaseModel

	UsedBy consts.UsedBy         `json:"usedBy"`
	Type   consts.CheckpointType `json:"type"`

	Expression        string `json:"expression"`
	ExtractorVariable string `json:"extractorVariable"`

	Operator consts.ComparisonOperator `json:"operator"`
	Value    string                    `json:"value"`

	ActualResult string              `json:"actualResult"`
	ResultStatus consts.ResultStatus `json:"resultStatus"`
	InterfaceId  uint                `json:"interfaceId"`
	ScenarioId   uint                `json:"scenarioId"`
}

func (DebugInterfaceCheckpoint) TableName() string {
	return "biz_debug_interface_checkpoint"
}