package service

import (
	serverDomain "github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	serverConsts "github.com/aaronchen2k/deeptest/internal/server/consts"
	model "github.com/aaronchen2k/deeptest/internal/server/modules/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/repo"
	"github.com/jinzhu/copier"
)

type DiagnoseInterfaceService struct {
	EndpointInterfaceRepo *repo.EndpointInterfaceRepo `inject:""`
	DiagnoseInterfaceRepo *repo.DiagnoseInterfaceRepo `inject:""`
	EndpointRepo          *repo.EndpointRepo          `inject:""`
	ServeRepo             *repo.ServeRepo             `inject:""`
	ScenarioProcessorRepo *repo.ScenarioProcessorRepo `inject:""`
	ServeServerRepo       *repo.ServeServerRepo       `inject:""`
	DebugInterfaceRepo    *repo.DebugInterfaceRepo    `inject:""`
	ExtractorRepo         *repo.ExtractorRepo         `inject:""`
	CheckpointRepo        *repo.CheckpointRepo        `inject:""`

	DebugInterfaceService *DebugInterfaceService `inject:""`
}

func (s *DiagnoseInterfaceService) Load(projectId, serveId int) (ret []*serverDomain.DiagnoseInterface, err error) {
	root, err := s.DiagnoseInterfaceRepo.GetTree(uint(projectId), uint(serveId))

	if root != nil {
		ret = root.Children
	}

	return
}

func (s *DiagnoseInterfaceService) Get(id int) (ret model.DiagnoseInterface, err error) {
	ret, err = s.DiagnoseInterfaceRepo.Get(uint(id))
	// its debug data will load in webpage

	return
}

func (s *DiagnoseInterfaceService) Save(req serverDomain.DiagnoseInterfaceSaveReq) (diagnoseInterface model.DiagnoseInterface, err error) {
	s.CopyValueFromRequest(&diagnoseInterface, req)

	if diagnoseInterface.Type == serverConsts.DiagnoseInterfaceTypeInterface {
		server, _ := s.ServeServerRepo.GetDefaultByServe(diagnoseInterface.ServeId)

		// create new DebugInterface
		debugInterface := model.DebugInterface{
			InterfaceBase: model.InterfaceBase{
				Name: req.Title,
				InterfaceConfigBase: model.InterfaceConfigBase{
					Method: consts.GET,
				},
			},
			ServeId:  diagnoseInterface.ServeId,
			ServerId: server.ID,
			BaseUrl:  server.Url,
		}

		err = s.DebugInterfaceRepo.Save(&debugInterface)
		diagnoseInterface.DebugInterfaceId = debugInterface.ID
	}

	err = s.DiagnoseInterfaceRepo.Save(&diagnoseInterface)

	if diagnoseInterface.DebugInterfaceId > 0 {
		values := map[string]interface{}{
			"diagnose_interface_id": diagnoseInterface.ID,
		}
		err = s.DebugInterfaceRepo.UpdateDebugInfo(diagnoseInterface.DebugInterfaceId, values)
	}

	return
}

func (s *DiagnoseInterfaceService) Remove(id int, typ serverConsts.DiagnoseInterfaceType) (err error) {
	err = s.DiagnoseInterfaceRepo.Remove(uint(id), typ)
	return
}

func (s *DiagnoseInterfaceService) Move(srcId, targetId uint, pos serverConsts.DropPos, projectId uint) (
	srcScenarioNode model.DiagnoseInterface, err error) {
	srcScenarioNode, err = s.DiagnoseInterfaceRepo.Get(srcId)

	srcScenarioNode.ParentId, srcScenarioNode.Ordr = s.DiagnoseInterfaceRepo.UpdateOrder(pos, targetId, projectId)
	err = s.DiagnoseInterfaceRepo.UpdateOrdAndParent(srcScenarioNode)

	return
}

func (s *DiagnoseInterfaceService) SaveDebugData(req domain.DebugData) (debugInterface model.DebugInterface, err error) {
	s.DebugInterfaceService.Save(req)

	return
}

func (s *DiagnoseInterfaceService) CopyValueFromRequest(interf *model.DiagnoseInterface, req serverDomain.DiagnoseInterfaceSaveReq) {
	copier.CopyWithOption(interf, req, copier.Option{
		DeepCopy: true,
	})
}

func (s *DiagnoseInterfaceService) CopyDebugDataValueFromRequest(interf *model.DiagnoseInterface, req domain.DebugData) (err error) {
	copier.CopyWithOption(interf, req, copier.Option{DeepCopy: true})

	return
}

func (s *DiagnoseInterfaceService) ImportInterfaces(req serverDomain.DiagnoseInterfaceImportReq) (ret model.DiagnoseInterface, err error) {
	parent, _ := s.DiagnoseInterfaceRepo.Get(req.TargetId)

	if parent.Type != serverConsts.DiagnoseInterfaceTypeDir {
		parent, _ = s.DiagnoseInterfaceRepo.Get(parent.ParentId)
	}

	for _, interfaceId := range req.InterfaceIds {
		ret, err = s.createInterfaceFromDefine(interfaceId, req.CreateBy, parent)
	}

	return
}

func (s *DiagnoseInterfaceService) createInterfaceFromDefine(endpointInterfaceId int, createBy uint, parent model.DiagnoseInterface) (
	ret model.DiagnoseInterface, err error) {

	endpointInterface, err := s.EndpointInterfaceRepo.Get(uint(endpointInterfaceId))
	if err != nil {
		return
	}

	// convert or clone a debug interface obj
	debugData, err := s.DebugInterfaceService.GetDebugDataFromEndpointInterface(uint(endpointInterfaceId))
	debugData.DebugInterfaceId = 0 // force to clone the old one
	debugData.DebugInterfaceId = 0
	debugData.EndpointInterfaceId = uint(endpointInterfaceId)
	debugData.ServeId = parent.ServeId

	server, _ := s.ServeServerRepo.GetDefaultByServe(debugData.ServeId)
	debugData.ServerId = server.ID
	debugData.BaseUrl = server.Url

	debugData.UsedBy = consts.DiagnoseDebug
	debugInterface, err := s.DebugInterfaceService.Save(debugData)

	// clone extractors and checkpoints if needed
	if endpointInterface.DebugInterfaceId <= 0 {
		s.ExtractorRepo.CloneFromEndpointInterfaceToDebugInterface(uint(endpointInterfaceId), debugInterface.ID, consts.DiagnoseDebug)
		s.CheckpointRepo.CloneFromEndpointInterfaceToDebugInterface(uint(endpointInterfaceId), debugInterface.ID, consts.DiagnoseDebug)
	}

	// save test interface
	diagnoseInterface := model.DiagnoseInterface{
		Title: endpointInterface.Name + "-" + string(endpointInterface.Method),
		Type:  serverConsts.DiagnoseInterfaceTypeInterface,
		Ordr:  s.DiagnoseInterfaceRepo.GetMaxOrder(parent.ID),

		DebugInterfaceId: debugInterface.ID,
		ParentId:         parent.ID,
		ServeId:          parent.ServeId,
		ProjectId:        parent.ProjectId,
		CreatedBy:        createBy,
	}
	s.DiagnoseInterfaceRepo.Save(&diagnoseInterface)

	// update diagnose_interface_id
	values := map[string]interface{}{
		"diagnose_interface_id": diagnoseInterface.ID,
	}
	s.DebugInterfaceRepo.UpdateDebugInfo(debugInterface.ID, values)

	ret = diagnoseInterface

	return
}