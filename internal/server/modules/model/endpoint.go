package model

type Endpoint struct {
	BaseModel
	Title       string              `json:"title"`
	ProjectId   int64               `json:"projectId"`
	ServeId     uint                `json:"serveId"`
	Path        string              `json:"path"`
	Version     string              `json:"version"`
	CreateUser  string              `json:"createUser"`
	Status      int64               `json:"status"`
	CategoryId  uint                `json:"categoryId"`
	PathParams  []EndpointPathParam `gorm:"-" json:"pathParams"`
	Interfaces  []EndpointInterface `gorm:"-" json:"interfaces"`
	Versions    []EndpointVersion   `gorm:"-" json:"versions"`
	ServeName   string              `gorm:"-" json:"serveName"`
	Description string              `json:"description"`
}

func (Endpoint) TableName() string {
	return "biz_endpoint"
}

type EndpointPathParam struct {
	BaseModel
	Name       string `json:"name"`
	Value      string `json:"value"`
	Type       string `json:"type"`
	EndpointId uint   `json:"endpointId"`
	Desc       string `json:"desc"`
	Required   bool   `json:"required"`
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
