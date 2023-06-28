package service

import (
	"fmt"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	repo "github.com/aaronchen2k/deeptest/internal/server/modules/repo"
)

type VariableService struct {
	DebugInterfaceRepo    *repo.DebugInterfaceRepo    `inject:""`
	DiagnoseInterfaceRepo *repo.DiagnoseInterfaceRepo `inject:""`

	EndpointInterfaceRepo *repo.EndpointInterfaceRepo `inject:""`
	EndpointRepo          *repo.EndpointRepo          `inject:""`

	ExtractorRepo   *repo.ExtractorRepo   `inject:""`
	EnvironmentRepo *repo.EnvironmentRepo `inject:""`
	ServeServerRepo *repo.ServeServerRepo `inject:""`

	EnvironmentService *EnvironmentService `inject:""`
	ShareVarService    *ShareVarService    `inject:""`
	DatapoolService    *DatapoolService    `inject:""`
}

func (s *VariableService) GetCombinedVarsForCheckpoint(debugInterfaceId, endpointInterfaceId, scenarioProcessorId uint, usedBy consts.UsedBy) (
	ret map[string]interface{}, datapools domain.Datapools, err error) {

	diagnoseInterfaceId := uint(0)

	if debugInterfaceId > 0 {
		debugInterface, _ := s.DebugInterfaceRepo.Get(debugInterfaceId)
		diagnoseInterfaceId = debugInterface.DiagnoseInterfaceId
	}

	server, _ := s.ServeServerRepo.GetByDebugInfo(debugInterfaceId, endpointInterfaceId)
	envId := server.EnvironmentId
	env, _ := s.EnvironmentRepo.Get(envId)
	projectId := env.ProjectId

	shareVariables := s.ShareVarService.List(debugInterfaceId, endpointInterfaceId, diagnoseInterfaceId, scenarioProcessorId, usedBy)
	envVars, _ := s.EnvironmentService.GetVarsByEnv(envId)
	globalVars, _ := s.EnvironmentService.GetGlobalVars(projectId)
	datapools, _ = s.DatapoolService.ListForExec(projectId)

	ret = CombineVariables(shareVariables, envVars, globalVars)

	return
}

func CombineVariables(shareVariables, envVars, globalVars []domain.GlobalVar) (
	ret map[string]interface{}) {
	ret = map[string]interface{}{}

	variableMap := map[string]interface{}{}

	for _, item := range globalVars {
		variableMap[item.Name] = item.LocalValue
	}
	for _, item := range envVars { // overwrite previous ones
		variableMap[item.Name] = item.LocalValue
	}
	for _, item := range shareVariables { // overwrite previous ones
		variableMap[item.Name] = item.LocalValue
	}

	// value is a  object
	for key, val := range variableMap {
		valMp, isMap := val.(map[string]interface{})

		if isMap {
			for propKey, v := range valMp {
				ret[fmt.Sprintf("%s.%s", key, propKey)] = v
			}

		} else {
			ret[fmt.Sprintf("%s", key)] = val

		}
	}

	return
}
