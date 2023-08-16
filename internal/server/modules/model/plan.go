package model

import "github.com/aaronchen2k/deeptest/internal/pkg/consts"

type Plan struct {
	BaseModel

	Version float64 `json:"version" yaml:"version"`
	Name    string  `json:"name" yaml:"name"`
	Desc    string  `json:"desc" yaml:"desc"`

	CategoryId     int               `json:"categoryId"`
	ProjectId      uint              `json:"projectId"`
	SerialNumber   string            `json:"serialNumber"`
	AdminId        uint              `json:"adminId"` //负责人ID
	CreateUserId   uint              `json:"createUserId"`
	UpdateUserId   uint              `json:"updateUserId"`
	Status         consts.TestStatus `json:"status"`
	TestStage      consts.TestStage  `json:"testStage"`
	Scenarios      []Scenario        `gorm:"-" json:"scenarios"`
	Reports        []PlanReport      `gorm:"-" json:"reports"`
	TestPassRate   string            `gorm:"-" json:"testPassRate"`
	AdminName      string            `gorm:"-" json:"adminName"`      //负责人姓名
	UpdateUserName string            `gorm:"-" json:"updateUserName"` //最近更新人姓名
	CurrEnvId      uint              `json:"currEnvId"`
	CreateUserName string            `gorm:"-" json:"createUserName"` //创建人姓名

}

func (Plan) TableName() string {
	return "biz_plan"
}

type RelaPlanScenario struct {
	BaseModel

	PlanId     uint `json:"planId"`
	ScenarioId uint `json:"scenarioId"`

	ServiceId uint `json:"serviceId"`
	ProjectId uint `json:"projectId"`
	SortId    uint `json:"sortId"` //排序ID
}

func (RelaPlanScenario) TableName() string {
	return "biz_plan_scenario_r"
}
