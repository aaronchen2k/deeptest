package model

import "github.com/aaronchen2k/deeptest/internal/pkg/consts"

type Endpoint struct {
	BaseModel
	Title        string              `json:"title"`
	ProjectId    uint                `json:"projectId" gorm:"index:idx_projectId_serveId"`
	ServeId      uint                `json:"serveId" gorm:"index:idx_projectId_serveId"`
	ServerId     uint                `json:"serverId"`
	Path         string              `json:"path"`
	Version      string              `json:"version"`
	CreateUser   string              `json:"createUser"`
	Status       int64               `json:"status"`
	CategoryId   int64               `json:"categoryId"`
	PathParams   []EndpointPathParam `gorm:"-" json:"pathParams"`
	Interfaces   []EndpointInterface `gorm:"-" json:"interfaces"`
	Versions     []EndpointVersion   `gorm:"-" json:"versions"`
	Tags         []string            `gorm:"-" json:"tags"`
	ServeName    string              `gorm:"-" json:"serveName"`
	Description  string              `json:"description"`
	SerialNumber string              `json:"serialNumber"`
	Curl         string              `gorm:"-" json:"curl"`
	SourceType   consts.SourceType   `json:"sourceType"`
	Maintainer   string              `gorm:"-" json:"maintainer"`
}

func (Endpoint) TableName() string {
	return "biz_endpoint"
}

type EndpointPathParam struct {
	EndpointInterfaceParam
	EndpointId uint `json:"endpointId"`
}

func (EndpointPathParam) TableName() string {
	return "biz_endpoint_path_param"
}

type EndpointVersion struct {
	BaseModel
	Version    string `json:"version"`
	EndpointId uint   `json:"endpointId"`
}

func (EndpointVersion) TableName() string {
	return "biz_endpoint_version"
}
