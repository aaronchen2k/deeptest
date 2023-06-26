package service

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/repo"
	"github.com/jinzhu/copier"
)

type ScenarioInterfaceService struct {
	EndpointInterfaceRepo *repo.EndpointInterfaceRepo `inject:""`
	DebugInterfaceRepo    *repo.DebugInterfaceRepo    `inject:""`
	ScenarioInterfaceRepo *repo.ScenarioInterfaceRepo `inject:""`
	EndpointRepo          *repo.EndpointRepo          `inject:""`
	ServeRepo             *repo.ServeRepo             `inject:""`
	ScenarioProcessorRepo *repo.ScenarioProcessorRepo `inject:""`
	ServeServerRepo       *repo.ServeServerRepo       `inject:""`
	TestInterfaceRepo     *repo.TestInterfaceRepo     `inject:""`

	ScenarioNodeService   *ScenarioNodeService   `inject:""`
	DebugSceneService     *DebugSceneService     `inject:""`
	DebugInterfaceService *DebugInterfaceService `inject:""`
	SceneService          *SceneService          `inject:""`
	EnvironmentService    *EnvironmentService    `inject:""`
	DatapoolService       *DatapoolService       `inject:""`
}

func (s *ScenarioInterfaceService) GetDebugDataFromScenarioInterface(scenarioInterfaceId uint) (req domain.DebugData, err error) {
	scenarioInterfacePo, _ := s.ScenarioInterfaceRepo.GetDetail(scenarioInterfaceId)
	if err != nil {
		return
	}

	endpointInterface, _ := s.EndpointInterfaceRepo.Get(scenarioInterfacePo.EndpointInterfaceId)

	s.SetProps(endpointInterface, &scenarioInterfacePo, &req)

	return
}

func (s *ScenarioInterfaceService) SetProps(
	endpointInterface model.EndpointInterface, scenarioInterfacePo *model.DebugInterface, debugData *domain.DebugData) {

	endpoint, err := s.EndpointRepo.GetAll(endpointInterface.EndpointId, "v0.1.0")
	serve, err := s.ServeRepo.Get(endpoint.ServeId)
	if err != nil {
		return
	}

	securities, err := s.ServeRepo.ListSecurity(serve.ID)
	if err != nil {
		return
	}

	serve.Securities = securities
	debugData.EndpointInterfaceId = endpointInterface.ID

	copier.CopyWithOption(&debugData, scenarioInterfacePo, copier.Option{DeepCopy: true})
	debugData.EndpointInterfaceId = endpointInterface.ID // reset

	debugData.Headers = append(debugData.Headers, domain.Header{Name: "", Value: ""})
	debugData.QueryParams = append(debugData.QueryParams, domain.Param{Name: "", Value: ""})
	debugData.PathParams = append(debugData.PathParams, domain.Param{Name: "", Value: ""})

	debugData.BodyFormData = append(debugData.BodyFormData, domain.BodyFormDataItem{
		Name: "", Value: "", Type: consts.FormDataTypeText})
	debugData.BodyFormUrlencoded = append(debugData.BodyFormUrlencoded, domain.BodyFormUrlEncodedItem{
		Name: "", Value: "",
	})

	return
}

//func (s *ScenarioInterfaceService) GetScenarioInterface(endpointInterfaceId uint) (ret domain.DebugData, err error) {
//	scenarioInterfaceId, _ := s.ScenarioInterfaceRepo.HasScenarioInterfaceRecord(endpointInterfaceId)
//
//	if scenarioInterfaceId > 0 {
//		ret, err = s.GetDebugDataFromScenarioInterface(scenarioInterfaceId)
//	} else {
//		ret, err = s.DebugInterfaceService.GetDebugInterfaceByEndpointInterface(endpointInterfaceId)
//		if err != nil || ret.EndpointInterfaceId == 0 {
//			return domain.DebugData{}, err
//		}
//		_, err = s.SaveDebugData(ret)
//	}
//
//	return
//}

func (s *ScenarioInterfaceService) SaveDebugData(req domain.DebugData) (debug model.DebugInterface, err error) {
	s.CopyValueFromRequest(&debug, req)

	endpointInterface, _ := s.EndpointInterfaceRepo.Get(req.EndpointInterfaceId)
	debug.EndpointId = endpointInterface.EndpointId

	if req.DebugInterfaceId > 0 {
		debug.ID = req.DebugInterfaceId
	}

	err = s.ScenarioInterfaceRepo.SaveDebugData(&debug)

	return
}

func (s *ScenarioInterfaceService) CopyValueFromRequest(interf *model.DebugInterface, req domain.DebugData) (err error) {
	copier.CopyWithOption(interf, req, copier.Option{DeepCopy: true})

	return
}

func (s *ScenarioInterfaceService) ResetDebugData(scenarioProcessorId int, createBy uint) (newProcessor model.Processor, err error) {
	scenarioProcessor, _ := s.ScenarioProcessorRepo.Get(uint(scenarioProcessorId))
	parentProcessor, _ := s.ScenarioProcessorRepo.Get(scenarioProcessor.ParentId)
	debugInterface, _ := s.DebugInterfaceRepo.Get(scenarioProcessor.EntityId)

	if debugInterface.TestInterfaceId > 0 {
		testInterface, _ := s.TestInterfaceRepo.Get(debugInterface.TestInterfaceId)
		testInterfaceTo := s.TestInterfaceRepo.ToTo(&testInterface)
		newProcessor, err = s.ScenarioNodeService.createDirOrInterfaceFromTest(testInterfaceTo, parentProcessor)

	} else if debugInterface.EndpointInterfaceId > 0 {
		serveId := uint(0)
		newProcessor, err = s.ScenarioNodeService.createInterfaceFromDefine(debugInterface.EndpointInterfaceId, &serveId, createBy, parentProcessor, scenarioProcessor.Name)
	}

	// must put below, since creation will use its DebugInterface
	s.DebugInterfaceRepo.Delete(scenarioProcessor.EntityId)
	s.ScenarioProcessorRepo.Delete(scenarioProcessor.ID)

	return
}
