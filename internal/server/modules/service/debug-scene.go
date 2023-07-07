package service

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/repo"
)

type DebugSceneService struct {
	EndpointInterfaceRepo *repo.EndpointInterfaceRepo `inject:""`
	EndpointRepo          *repo.EndpointRepo          `inject:""`
	ServeServerRepo       *repo.ServeServerRepo       `inject:""`
	ScenarioProcessorRepo *repo.ScenarioProcessorRepo `inject:""`
	EnvironmentRepo       *repo.EnvironmentRepo       `inject:""`
	DiagnoseInterfaceRepo *repo.DiagnoseInterfaceRepo `inject:""`

	ShareVarService *ShareVarService `inject:""`

	EnvironmentService *EnvironmentService `inject:""`
}

func (s *DebugSceneService) LoadScene(debugData *domain.DebugData) (
	baseUrl string, shareVars []domain.GlobalVar, envVars []domain.GlobalVar,
	globalVars []domain.GlobalVar, globalParams []domain.GlobalParam) {

	debugServeId := debugData.ServeId
	debugServerId := debugData.ServerId

	if debugData.EndpointInterfaceId > 0 && (debugServeId <= 0 || debugServerId <= 0) {
		interf, _ := s.EndpointInterfaceRepo.Get(debugData.EndpointInterfaceId)
		endpoint, _ := s.EndpointRepo.Get(interf.EndpointId)

		if debugServeId <= 0 {
			debugServeId = endpoint.ServeId
		}
		if debugServerId <= 0 {
			debugServerId = endpoint.ServerId
		}
	}

	serveServer, _ := s.ServeServerRepo.Get(debugServerId)

	if debugData.DiagnoseInterfaceId > 0 {
		baseUrl = debugData.BaseUrl
	} else {
		baseUrl = serveServer.Url
	}

	envId := serveServer.EnvironmentId
	environment, _ := s.EnvironmentRepo.Get(envId)
	debugData.ProjectId = environment.ProjectId

	shareVars, _ = s.ShareVarService.ListForDebug(debugServeId, debugData.ScenarioProcessorId, debugData.UsedBy)
	envVars, _ = s.EnvironmentService.GetVarsByEnv(envId)
	globalVars, _ = s.EnvironmentService.GetGlobalVars(environment.ProjectId)
	globalParams, _ = s.EnvironmentService.GetGlobalParams(environment.ProjectId)

	return
}
