package service

import (
	serverDomain "github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	model "github.com/aaronchen2k/deeptest/internal/server/modules/model"
	repo "github.com/aaronchen2k/deeptest/internal/server/modules/repo"
)

type PostConditionService struct {
	PostConditionRepo *repo.PostConditionRepo `inject:""`
	ExtractorRepo     *repo.ExtractorRepo     `inject:""`
	CheckpointRepo    *repo.CheckpointRepo    `inject:""`
	ScriptRepo        *repo.ScriptRepo        `inject:""`
}

func (s *PostConditionService) List(debugInterfaceId, endpointInterfaceId uint, category consts.ConditionCategory) (conditions []model.DebugPostCondition, err error) {
	conditions, err = s.PostConditionRepo.List(debugInterfaceId, endpointInterfaceId, category)

	return
}

func (s *PostConditionService) Get(id uint) (condition model.DebugPostCondition, err error) {
	condition, err = s.PostConditionRepo.Get(id)

	return
}

func (s *PostConditionService) Create(condition *model.DebugPostCondition) (err error) {
	err = s.PostConditionRepo.Save(condition)

	var entityId uint

	if condition.EntityType == consts.ConditionTypeExtractor {
		po := s.ExtractorRepo.CreateDefault(condition.ID)
		entityId = po.ID

	} else if condition.EntityType == consts.ConditionTypeCheckpoint {
		po := s.CheckpointRepo.CreateDefault(condition.ID)
		entityId = po.ID

	} else if condition.EntityType == consts.ConditionTypeScript {
		po := s.ScriptRepo.CreateDefault(condition.ID, consts.ConditionSrcPost)
		entityId = po.ID
	}

	err = s.PostConditionRepo.UpdateEntityId(condition.ID, entityId)

	return
}

//func (s *PostConditionService) CloneAll(srcDebugInterfaceId, srcEndpointInterfaceId, distDebugInterfaceId uint) (err error) {
//	return s.PostConditionRepo.CloneAll(srcDebugInterfaceId, srcEndpointInterfaceId, distDebugInterfaceId)
//}

func (s *PostConditionService) Delete(reqId uint) (err error) {
	err = s.PostConditionRepo.Delete(reqId)

	return
}

func (s *PostConditionService) Disable(reqId uint) (err error) {
	err = s.PostConditionRepo.Disable(reqId)

	return
}

func (s *PostConditionService) Move(req serverDomain.ConditionMoveReq) (err error) {
	err = s.PostConditionRepo.UpdateOrders(req)

	return
}