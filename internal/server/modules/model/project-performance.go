package model

type ProjectPerformanceSetting struct {
	BaseModel

	ConductorGrpcAddress string `json:"conductorGrpcAddress"`

	InfluxdbAddress string `json:"influxdbAddress"`
	InfluxdbOrg     string `json:"influxdbOrg"`
	InfluxdbToken   string `json:"influxdbToken"`

	ProjectId uint `json:"projectId"`
}

func (ProjectPerformanceSetting) TableName() string {
	return "biz_project_performance_settings"
}
