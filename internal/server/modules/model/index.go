package model

import (
	"github.com/aaronchen2k/deeptest/internal/server/middleware"
)

var (
	Models = []interface{}{
		&middleware.Oplog{},

		&SysPerm{},
		&SysRole{},
		&SysUser{},
		&SysUserProfile{},

		&ProjectRole{},
		&Org{},
		&Project{},
		&ProjectMember{},
		&Datapool{},
		&Environment{},
		&EnvironmentVar{},
		&ShareVariable{},

		&DebugInterface{},
		&DebugInterfaceParam{},
		&DebugInterfaceBodyFormDataItem{},
		&DebugInterfaceBodyFormUrlEncodedItem{},
		&DebugInterfaceHeader{},
		&DebugInterfaceBasicAuth{},
		&DebugInterfaceBearerToken{},
		&DebugInterfaceOAuth20{},
		&DebugInterfaceApiKey{},

		&DebugPreCondition{},
		&DebugPostCondition{},
		&DebugConditionExtractor{},
		&DebugConditionCheckpoint{},
		&DebugConditionScript{},

		&DiagnoseInterface{},

		&Snippet{},

		&MockInvocation{},
		&Auth2Token{},

		&Category{},
		&Scenario{},

		&Plan{},
		&RelaPlanScenario{},

		&Processor{},
		//&ProcessorThreadGroup{},
		&ProcessorGroup{},
		&ProcessorLogic{},
		&ProcessorLoop{},
		&ProcessorTimer{},
		&ProcessorPrint{},
		&ProcessorVariable{},
		&ProcessorAssertion{},
		&ProcessorData{},
		&ProcessorCookie{},
		&ProcessorExtractor{},

		&ScenarioReport{},
		&PlanReport{},
		&ExecLogProcessor{},
		&ExecLogExtractor{},
		&ExecLogCheckpoint{},

		&ComponentSchema{},
		&ComponentSchemaSecurity{},

		&Endpoint{},
		&EndpointPathParam{},
		&EndpointInterfaceRequestBody{},
		&EndpointInterfaceRequestBodyItem{},
		&EndpointInterfaceResponseBodyItem{},
		&EndpointInterfaceResponseBodyHeader{},
		&EndpointInterfaceResponseBody{},
		&EndpointInterface{},
		&EndpointCase{},
		&EndpointInterfaceParam{},
		&EndpointInterfaceCookie{},
		&EndpointInterfaceHeader{},
		&EndpointDocument{},
		&EndpointSnapshot{},

		&Serve{},
		&ServeServer{},
		&ServeVersion{},
		&EndpointVersion{},
		&ServeEndpointVersion{},
		&SummaryBugs{},
		&SummaryDetails{},
		&SummaryProjectUserRanking{},
		&EnvironmentParam{},
		&Message{},
		&MessageRead{},
		&DebugInvoke{},
		&ProjectPerm{},
		&ProjectRolePerm{},
		&ProjectRoleMenu{},
		&ProjectMenu{},
		&ProjectRecentlyVisited{},
		&ProjectMemberAudit{},

		&SwaggerSync{},
	}
)
