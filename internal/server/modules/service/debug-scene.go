package service

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/repo"
	_httpUtils "github.com/aaronchen2k/deeptest/pkg/lib/http"
)

type DebugSceneService struct {
	EndpointInterfaceRepo *repo.EndpointInterfaceRepo `inject:""`
	EndpointRepo          *repo.EndpointRepo          `inject:""`
	ServeServerRepo       *repo.ServeServerRepo       `inject:""`
	ScenarioProcessorRepo *repo.ScenarioProcessorRepo `inject:""`
	EnvironmentRepo       *repo.EnvironmentRepo       `inject:""`

	ShareVarService *ShareVarService `inject:""`

	EnvironmentService *EnvironmentService `inject:""`
}

func (s *DebugSceneService) LoadScene(interfaceId, endpointId, processorId uint, usedBy consts.UsedBy) (
	baseUrl string, shareVariables []domain.ShareVars, envVars []domain.EnvVars,
	globalEnvVars []domain.GlobalEnvVars, globalParamVars []domain.GlobalParamVars) {

	var serveId, serverId, scenarioId, projectId uint

	interf, _ := s.EndpointInterfaceRepo.Get(interfaceId)
	endpoint, _ := s.EndpointRepo.Get(interf.EndpointId)
	serveId = endpoint.ServeId
	serverId = endpoint.ServerId
	projectId = endpoint.ProjectId

	serveServer, _ := s.ServeServerRepo.Get(serverId)
	baseUrl = _httpUtils.AddSepIfNeeded(serveServer.Url)
	envId := serveServer.EnvironmentId

	// by scenario
	if usedBy == consts.ScenarioDebug {
		processor, _ := s.ScenarioProcessorRepo.Get(processorId)
		scenarioId = processor.ScenarioId
	}

	shareVariables, _ = s.ShareVarService.ListForDebug(serveId, scenarioId, usedBy)
	envVars, _ = s.EnvironmentService.GetVarsByEnv(envId)

	globalEnvVars, _ = s.EnvironmentService.GetGlobalVars(projectId)
	globalParamVars, _ = s.EnvironmentService.GetGlobalParams(projectId)

	// interf, _ := s.ProcessorInterfaceRepo.Get(req.InterfaceId)
	//req.Datapools, _ = s.DatapoolService.ListForExec(interf.ProjectId)

	return
}
