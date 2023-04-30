package service

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/model"
	repo "github.com/aaronchen2k/deeptest/internal/server/modules/repo"
)

type ShareVarService struct {
	ShareVariableRepo *repo.ShareVariableRepo `inject:""`

	EndpointInterfaceRepo *repo.EndpointInterfaceRepo `inject:""`
	EndpointRepo          *repo.EndpointRepo          `inject:""`
	ServeServerRepo       *repo.ServeServerRepo       `inject:""`
	ScenarioProcessorRepo *repo.ScenarioProcessorRepo `inject:""`
}

func (s *ShareVarService) Save(name, value string, interfaceId, serveId, processorId, scenarioId uint,
	scope consts.ExtractorScope, usedBy consts.UsedBy) (err error) {

	po := model.ShareVariable{
		Name:        name,
		Value:       value,
		InterfaceId: interfaceId,
		ServeId:     serveId,
		ProcessorId: processorId,
		ScenarioId:  scenarioId,
		Scope:       scope,
	}

	if usedBy == consts.InterfaceDebug {
		po.ID, err = s.ShareVariableRepo.GetExistByInterfaceDebug(name, serveId)
	} else if usedBy == consts.ScenarioDebug {
		po.ID, err = s.ShareVariableRepo.GetExistByScenarioDebug(name, scenarioId)
	}

	err = s.ShareVariableRepo.Save(&po)

	return
}

func (s *ShareVarService) List(endpointInterfaceId, scenarioProcessorId uint, usedBy consts.UsedBy) (
	shareVariables []domain.VarKeyValuePair) {

	interf, _ := s.EndpointInterfaceRepo.Get(endpointInterfaceId)
	endpoint, _ := s.EndpointRepo.Get(interf.EndpointId)
	serveId := endpoint.ServeId

	shareVariables, _ = s.listForDebug(serveId, scenarioProcessorId, usedBy)

	return
}

func (s *ShareVarService) Delete(id int) (err error) {
	err = s.ShareVariableRepo.Delete(id)

	return
}

func (s *ShareVarService) Clear(endpointOrProcessorId int, usedBy consts.UsedBy) (err error) {
	if usedBy == consts.InterfaceDebug {
		endpoint, _ := s.EndpointRepo.Get(uint(endpointOrProcessorId))
		err = s.ShareVariableRepo.DeleteAllByServeId(endpoint.ServeId)

	} else if usedBy == consts.ScenarioDebug {
		processor, _ := s.ScenarioProcessorRepo.Get(uint(endpointOrProcessorId))
		err = s.ShareVariableRepo.DeleteAllByScenarioId(processor.ScenarioId)

	}

	return
}

func (s *ShareVarService) listForDebug(serveId, scenarioProcessorId uint, usedBy consts.UsedBy) (ret []domain.VarKeyValuePair, err error) {
	var pos []model.ShareVariable

	if usedBy == consts.InterfaceDebug {
		pos, err = s.ShareVariableRepo.ListByInterfaceDebug(serveId)
	} else if usedBy == consts.ScenarioDebug {
		pos, err = s.ShareVariableRepo.ListByScenarioDebug(scenarioProcessorId)
	}

	for _, po := range pos {
		ret = append(ret, domain.VarKeyValuePair{
			"id":    po.ID,
			"name":  po.Name,
			"value": po.Value,
		})
	}

	return
}
