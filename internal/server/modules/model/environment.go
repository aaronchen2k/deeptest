package model

import "github.com/aaronchen2k/deeptest/internal/pkg/domain"

type Environment struct {
	BaseModel
	ProjectId    uint             `json:"projectId"`
	Name         string           `json:"name"`
	Vars         []EnvironmentVar `gorm:"-" json:"vars"`
	ServeServers []ServeServer    `gorm:"-" json:"serveServers"`
	Sort         uint             `json:"sort"`
}

type EnvironmentVar struct {
	BaseModel
	domain.GlobalVar

	EnvironmentId uint `json:"environmentId"`
	ProjectId     uint `json:"projectId"`
}

type EnvironmentParam struct {
	BaseModel
	domain.GlobalParam

	Description string `json:"description"`
	ProjectId   uint   `json:"projectId"`
}

func (Environment) TableName() string {
	return "biz_environment"
}
func (EnvironmentVar) TableName() string {
	return "biz_environment_var"
}

func (EnvironmentParam) TableName() string {
	return "biz_environment_param"
}
