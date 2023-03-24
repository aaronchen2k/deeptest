package agentExec

type PlanExecReq struct {
	PlanId uint `json:"planId"`

	ServerUrl string `json:"serverUrl"`
	Token     string `json:"token"`
}

type PlanExecObj struct {
	Name      string            `json:"name"`
	Scenarios []ScenarioExecObj `json:"scenarios"`

	ServerUrl string `json:"serverUrl"`
	Token     string `json:"token"`
}