package service

import (
	"fmt"
	agentExec "github.com/aaronchen2k/deeptest/internal/agent/exec"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/helper/openapi"
	model "github.com/aaronchen2k/deeptest/internal/server/modules/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/repo"
	"github.com/jinzhu/copier"
)

type DebugInterfaceService struct {
	EndpointInterfaceRepo *repo.EndpointInterfaceRepo `inject:""`
	DebugInterfaceRepo    *repo.DebugInterfaceRepo    `inject:""`
	ScenarioInterfaceRepo *repo.ScenarioInterfaceRepo `inject:""`
	DiagnoseInterfaceRepo *repo.DiagnoseInterfaceRepo `inject:""`
	EndpointCaseRepo      *repo.EndpointCaseRepo      `inject:""`

	EndpointRepo          *repo.EndpointRepo          `inject:""`
	ServeRepo             *repo.ServeRepo             `inject:""`
	ScenarioProcessorRepo *repo.ScenarioProcessorRepo `inject:""`
	ServeServerRepo       *repo.ServeServerRepo       `inject:""`
	ExtractorRepo         *repo.ExtractorRepo         `inject:""`
	CheckpointRepo        *repo.CheckpointRepo        `inject:""`

	DebugSceneService        *DebugSceneService        `inject:""`
	ScenarioInterfaceService *ScenarioInterfaceService `inject:""`
	SceneService             *SceneService             `inject:""`
	EnvironmentService       *EnvironmentService       `inject:""`
	DatapoolService          *DatapoolService          `inject:""`
	ServeService             *ServeService             `inject:""`
}

func (s *DebugInterfaceService) Load(loadReq domain.DebugReq) (debugData domain.DebugData, err error) {
	if loadReq.DebugInterfaceId > 0 {
		debugData, _ = s.GetDebugDataFromDebugInterface(loadReq.DebugInterfaceId)
	} else {
		if loadReq.ScenarioProcessorId > 0 {
			debugData, _ = s.GetDebugInterfaceByScenarioInterface(loadReq.ScenarioProcessorId)
		} else if loadReq.DiagnoseInterfaceId > 0 {
			debugData, _ = s.GetDebugInterfaceByDiagnoseInterface(loadReq.DiagnoseInterfaceId)
		} else if loadReq.CaseInterfaceId > 0 {
			debugData, _ = s.GetDebugInterfaceByEndpointCase(loadReq.CaseInterfaceId)
		} else if loadReq.EndpointInterfaceId > 0 {
			debugData, _ = s.GetDebugInterfaceByEndpointInterface(loadReq.EndpointInterfaceId)
		}
	}

	debugData.UsedBy = loadReq.UsedBy

	debugData.BaseUrl, debugData.ShareVars, debugData.EnvVars, debugData.GlobalVars, debugData.GlobalParams =
		s.DebugSceneService.LoadScene(&debugData)

	return
}

func (s *DebugInterfaceService) LoadForExec(loadReq domain.DebugReq) (ret agentExec.InterfaceExecObj, err error) {
	ret.DebugData, _ = s.Load(loadReq)

	ret.ExecScene.ShareVars = ret.DebugData.ShareVars // for execution
	ret.DebugData.ShareVars = nil                     // for display on debug page only

	// get project environment data
	s.SceneService.LoadEnvVars(&ret.ExecScene, ret.DebugData)

	s.SceneService.LoadProjectSettings(&ret.ExecScene, ret.DebugData.ProjectId)

	return
}

func (s *DebugInterfaceService) GetDetail(id uint) (ret model.DebugInterface, err error) {
	ret, err = s.DebugInterfaceRepo.GetDetail(id)

	return
}

func (s *DebugInterfaceService) Save(req domain.DebugData) (debugInterface model.DebugInterface, err error) {
	s.CopyValueFromRequest(&debugInterface, req)

	if req.DebugInterfaceId > 0 { // to update
		debugInterface.ID = req.DebugInterfaceId
	} else {
		debugInterface.ServeId = req.ServeId
	}

	err = s.DebugInterfaceRepo.Save(&debugInterface)

	if req.DebugInterfaceId <= 0 { // first time to save
		// clone extractors and checkpoints if needed
		s.ExtractorRepo.CloneFromEndpointInterfaceToDebugInterface(req.EndpointInterfaceId, debugInterface.ID, req.UsedBy)
		s.CheckpointRepo.CloneFromEndpointInterfaceToDebugInterface(req.EndpointInterfaceId, debugInterface.ID, req.UsedBy)
	}

	if req.UsedBy == consts.InterfaceDebug {
		s.EndpointInterfaceRepo.SetDebugInterfaceId(req.EndpointInterfaceId, debugInterface.ID)
	}

	return
}

func (s *DebugInterfaceService) GetDebugDataFromDebugInterface(debugInterfaceId uint) (debugData domain.DebugData, err error) {
	debugInterfacePo, err := s.DebugInterfaceRepo.GetDetail(debugInterfaceId)
	if err != nil {
		return
	}

	endpointInterface, _ := s.EndpointInterfaceRepo.Get(debugInterfacePo.EndpointInterfaceId)

	s.CopyDebugDataPropsFromPo(&debugData, &debugInterfacePo, endpointInterface)

	return
}
func (s *DebugInterfaceService) GetDebugDataFromEndpointInterface(endpointInterfaceId uint) (ret domain.DebugData, err error) {
	endpointInterface, _ := s.EndpointInterfaceRepo.Get(endpointInterfaceId)

	if endpointInterface.DebugInterfaceId > 0 {
		ret, err = s.GetDebugDataFromDebugInterface(endpointInterface.DebugInterfaceId)
	} else {
		ret, err = s.ConvertDebugDataFromEndpointInterface(endpointInterfaceId)
	}

	return
}

func (s *DebugInterfaceService) ConvertDebugDataFromEndpointInterface(endpointInterfaceId uint) (debugData domain.DebugData, err error) {
	endpointInterface, err := s.EndpointInterfaceRepo.GetDetail(endpointInterfaceId)
	if err != nil {
		return
	}

	s.CopyDebugDataPropsFromPo(&debugData, nil, endpointInterface)

	// init ServeId and ServerId id by endpoint
	endpoint, _ := s.EndpointRepo.Get(endpointInterface.EndpointId)
	debugData.ServeId = endpoint.ServeId
	debugData.ServerId = endpoint.ServerId
	debugData.ProjectId = endpoint.ProjectId

	debugData.UsedBy = consts.InterfaceDebug

	return
}

func (s *DebugInterfaceService) CopyDebugDataPropsFromPo(debugData *domain.DebugData,
	debugInterfacePo *model.DebugInterface, endpointInterface model.EndpointInterface) {

	endpoint, err := s.EndpointRepo.GetAll(endpointInterface.EndpointId, "v0.1.0")
	serve, err := s.ServeRepo.Get(endpoint.ServeId)
	if err != nil {
		//	return
	}

	securities, err := s.ServeRepo.ListSecurity(serve.ID)
	if err != nil {
		//	return
	}

	serve.Securities = securities
	debugData.EndpointInterfaceId = endpointInterface.ID

	if debugInterfacePo == nil { // is null when converting from EndpointInterface
		schema2conv := openapi.NewSchema2conv()
		schema2conv.Components = s.ServeService.Components(serve.ID)
		interfaces2debug := openapi.NewInterfaces2debug(endpointInterface, endpoint, serve, schema2conv)
		debugInterfacePo = interfaces2debug.Convert()

		debugInterfacePo.Name = fmt.Sprintf("%s - %s", endpoint.Title, debugInterfacePo.Method)
	}

	copier.CopyWithOption(&debugData, debugInterfacePo, copier.Option{DeepCopy: true})
	debugData.EndpointInterfaceId = endpointInterface.ID // reset

	debugData.DebugInterfaceId = debugInterfacePo.ID
	debugData.ServeId = debugInterfacePo.ServeId

	debugData.Headers = append(debugData.Headers, domain.Header{Name: "", Value: ""})
	debugData.QueryParams = append(debugData.QueryParams, domain.Param{Name: "", Value: "", ParamIn: consts.ParamInQuery})
	debugData.PathParams = append(debugData.PathParams, domain.Param{Name: "", Value: "", ParamIn: consts.ParamInPath})

	debugData.BodyFormData = append(debugData.BodyFormData, domain.BodyFormDataItem{
		Name: "", Value: "", Type: consts.FormDataTypeText})
	debugData.BodyFormUrlencoded = append(debugData.BodyFormUrlencoded, domain.BodyFormUrlEncodedItem{
		Name: "", Value: "",
	})

	return
}

func (s *DebugInterfaceService) GetEndpointAndServeIdForEndpointInterface(endpointInterfaceId uint) (
	endpointId, serveId uint) {

	endpointInterface, _ := s.EndpointInterfaceRepo.Get(endpointInterfaceId)

	endpointId = endpointInterface.EndpointId
	endpoint, _ := s.EndpointRepo.Get(endpointId)

	serveId = endpoint.ServeId

	return
}

func (s *DebugInterfaceService) GetScenarioIdForDebugInterface(processorId uint) (
	scenarioId uint) {

	processor, _ := s.ScenarioProcessorRepo.Get(processorId)
	scenarioId = processor.ScenarioId

	return
}

func (s *DebugInterfaceService) CopyValueFromRequest(interf *model.DebugInterface, req domain.DebugData) (err error) {
	copier.CopyWithOption(interf, req, copier.Option{DeepCopy: true})

	return
}

func (s *DebugInterfaceService) GenSample(projectId, serveId uint) (ret *model.DebugInterface, err error) {
	return
}

func (s *DebugInterfaceService) GetDebugInterfaceByEndpointInterface(endpointInterfaceId uint) (ret domain.DebugData, err error) {
	endpointInterface, _ := s.EndpointInterfaceRepo.Get(endpointInterfaceId)

	if endpointInterface.DebugInterfaceId > 0 {
		ret, err = s.GetDebugDataFromDebugInterface(endpointInterface.DebugInterfaceId)
	} else {
		ret, err = s.ConvertDebugDataFromEndpointInterface(endpointInterfaceId)
	}

	return
}

func (s *DebugInterfaceService) GetDebugInterfaceByScenarioInterface(scenarioProcessorId uint) (ret domain.DebugData, err error) {
	processor, err := s.ScenarioProcessorRepo.Get(scenarioProcessorId)
	if err != nil {
		return
	}

	debugInterfaceId := processor.EntityId
	debugData, _ := s.DebugInterfaceRepo.GetDetail(debugInterfaceId)

	copier.CopyWithOption(&ret, debugData, copier.Option{
		DeepCopy: true,
	})

	if ret.ServerId <= 0 {
		endpointInterface, _ := s.EndpointInterfaceRepo.Get(processor.EndpointInterfaceId)
		server, _ := s.ServeServerRepo.GetByEndpoint(endpointInterface.EndpointId)
		ret.ServerId = server.ID
	}

	ret.DebugInterfaceId = debugInterfaceId
	ret.ScenarioProcessorId = scenarioProcessorId
	ret.ProcessorInterfaceSrc = debugData.ProcessorInterfaceSrc

	ret.Headers = append(ret.Headers, domain.Header{Name: "", Value: ""})
	ret.QueryParams = append(ret.QueryParams, domain.Param{Name: "", Value: "", ParamIn: consts.ParamInQuery})
	ret.PathParams = append(ret.PathParams, domain.Param{Name: "", Value: "", ParamIn: consts.ParamInPath})

	ret.BodyFormData = append(ret.BodyFormData, domain.BodyFormDataItem{
		Name: "", Value: "", Type: consts.FormDataTypeText})
	ret.BodyFormUrlencoded = append(ret.BodyFormUrlencoded, domain.BodyFormUrlEncodedItem{
		Name: "", Value: "",
	})

	return
}

func (s *DebugInterfaceService) GetDebugInterfaceByDiagnoseInterface(diagnoseInterfaceId uint) (ret domain.DebugData, err error) {
	diagnoseInterface, err := s.DiagnoseInterfaceRepo.GetDetail(diagnoseInterfaceId)
	if err != nil {
		return
	}

	copier.CopyWithOption(&ret, diagnoseInterface.DebugData, copier.Option{
		DeepCopy: true,
	})

	ret.ServerId = diagnoseInterface.DebugData.ServerId
	ret.DebugInterfaceId = diagnoseInterface.DebugInterfaceId
	ret.DiagnoseInterfaceId = diagnoseInterfaceId

	ret.Headers = append(ret.Headers, domain.Header{Name: "", Value: ""})
	ret.QueryParams = append(ret.QueryParams, domain.Param{Name: "", Value: "", ParamIn: consts.ParamInQuery})
	ret.PathParams = append(ret.PathParams, domain.Param{Name: "", Value: "", ParamIn: consts.ParamInPath})

	ret.BodyFormData = append(ret.BodyFormData, domain.BodyFormDataItem{
		Name: "", Value: "", Type: consts.FormDataTypeText})
	ret.BodyFormUrlencoded = append(ret.BodyFormUrlencoded, domain.BodyFormUrlEncodedItem{
		Name: "", Value: "",
	})

	return
}

func (s *DebugInterfaceService) GetDebugInterfaceByEndpointCase(endpointCaseId uint) (ret domain.DebugData, err error) {
	endpointCase, err := s.EndpointCaseRepo.GetDetail(endpointCaseId)
	if err != nil {
		return
	}

	copier.CopyWithOption(&ret, endpointCase.DebugData, copier.Option{
		DeepCopy: true,
	})

	ret.ServerId = endpointCase.DebugData.ServerId
	ret.DebugInterfaceId = endpointCase.DebugInterfaceId
	ret.CaseInterfaceId = endpointCaseId

	ret.Headers = append(ret.Headers, domain.Header{Name: "", Value: ""})
	ret.QueryParams = append(ret.QueryParams, domain.Param{Name: "", Value: "", ParamIn: consts.ParamInQuery})
	ret.PathParams = append(ret.PathParams, domain.Param{Name: "", Value: "", ParamIn: consts.ParamInPath})

	ret.BodyFormData = append(ret.BodyFormData, domain.BodyFormDataItem{
		Name: "", Value: "", Type: consts.FormDataTypeText})
	ret.BodyFormUrlencoded = append(ret.BodyFormUrlencoded, domain.BodyFormUrlEncodedItem{
		Name: "", Value: "",
	})

	return
}
