package service

import (
	serverDomain "github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	model "github.com/aaronchen2k/deeptest/internal/server/modules/model"
	repo "github.com/aaronchen2k/deeptest/internal/server/modules/repo"
)

type PreConditionService struct {
	PreConditionRepo *repo.PreConditionRepo `inject:""`
	ScriptRepo       *repo.ScriptRepo       `inject:""`
}

func (s *PreConditionService) GetScript(debugInterfaceId, endpointInterfaceId uint) (script model.DebugConditionScript, err error) {
	conditions, err := s.PreConditionRepo.List(debugInterfaceId, endpointInterfaceId)

	if len(conditions) == 0 {
		condition := model.DebugPreCondition{
			DebugInterfaceId:    debugInterfaceId,
			EndpointInterfaceId: endpointInterfaceId,
			EntityType:          consts.ConditionTypeScript,
		}
		err = s.Create(&condition)

		script, _ = s.ScriptRepo.GetByCondition(condition.ID)
		if script.ID == 0 {
			script = s.ScriptRepo.CreateDefault(condition.ID, consts.ConditionSrcPre)
		}

		s.PreConditionRepo.UpdateEntityId(condition.ID, script.ID)

		conditions, err = s.PreConditionRepo.List(debugInterfaceId, endpointInterfaceId)
	}

	script, err = s.ScriptRepo.Get(conditions[0].EntityId)
	script.ConditionEntityType = consts.ConditionTypeScript

	return
}

func (s *PreConditionService) Get(id uint) (checkpoint model.DebugPreCondition, err error) {
	checkpoint, err = s.PreConditionRepo.Get(id)

	return
}

func (s *PreConditionService) Create(condition *model.DebugPreCondition) (err error) {
	err = s.PreConditionRepo.Save(condition)

	var entityId uint

	if condition.EntityType == consts.ConditionTypeScript {
		po := s.ScriptRepo.CreateDefault(condition.ID, consts.ConditionSrcPre)
		entityId = po.ID
	}

	err = s.PreConditionRepo.UpdateEntityId(condition.ID, entityId)

	return
}

//func (s *PreConditionService) CloneAll(srcDebugInterfaceId, srcEndpointInterfaceId, distDebugInterfaceId uint) (err error) {
//	return s.PreConditionRepo.CloneAll(srcDebugInterfaceId, srcEndpointInterfaceId, distDebugInterfaceId)
//}

func (s *PreConditionService) Delete(reqId uint) (err error) {
	err = s.PreConditionRepo.Delete(reqId)

	return
}

func (s *PreConditionService) Disable(reqId uint) (err error) {
	err = s.PreConditionRepo.Disable(reqId)

	return
}

func (s *PreConditionService) Move(req serverDomain.ConditionMoveReq) (err error) {
	err = s.PreConditionRepo.UpdateOrders(req)

	return
}
