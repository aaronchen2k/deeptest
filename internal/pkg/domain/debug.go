package domain

import "github.com/aaronchen2k/deeptest/internal/pkg/consts"

type DebugReq struct {
	EndpointInterfaceId uint `json:"endpointInterfaceId"` // load by endpoint designer
	ScenarioProcessorId uint `json:"scenarioProcessorId"` // load by scenario designer

	UsedBy consts.UsedBy `json:"usedBy"`
}

type SubmitDebugResultRequest struct {
	Request  DebugData     `json:"request"`
	Response DebugResponse `json:"response"`
}

type DebugData struct {
	EndpointInterfaceId uint          `gorm:"-" json:"endpointInterfaceId"`
	ScenarioProcessorId uint          `gorm:"-" json:"scenarioProcessorId"`
	UsedBy              consts.UsedBy `json:"usedBy"`

	BaseUrl   string            `gorm:"-" json:"baseUrl"`
	ShareVars []VarKeyValuePair `gorm:"-" json:"shareVars"`

	//EnvVars         []VarKeyValuePair `gorm:"-" json:"envVars"`
	//GlobalEnvVars   []GlobalVars       `gorm:"-" json:"globalEnvVars"`
	//GlobalParamVars []GlobalParams     `gorm:"-" json:"globalParamVars"`
	//Datapools   Datapools   `gorm:"-" json:"datapools"`

	Name string `gorm:"-" json:"name"`
	BaseRequest
}
