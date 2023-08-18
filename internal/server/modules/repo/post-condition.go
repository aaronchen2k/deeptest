package repo

import (
	"encoding/json"
	"fmt"
	serverDomain "github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	model "github.com/aaronchen2k/deeptest/internal/server/modules/model"
	commonUtils "github.com/aaronchen2k/deeptest/pkg/lib/comm"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type PostConditionRepo struct {
	DB *gorm.DB `inject:""`

	ExtractorRepo         *ExtractorRepo         `inject:""`
	CookieRepo            *CookieRepo            `inject:""`
	CheckpointRepo        *CheckpointRepo        `inject:""`
	ScriptRepo            *ScriptRepo            `inject:""`
	ResponseDefineRepo    *ResponseDefineRepo    `inject:""`
	EndpointInterfaceRepo *EndpointInterfaceRepo `inject:""`
}

func (r *PostConditionRepo) List(debugInterfaceId, endpointInterfaceId uint, typ consts.ConditionCategory) (pos []model.DebugPostCondition, err error) {
	db := r.DB.Where("NOT deleted")

	if debugInterfaceId > 0 {
		db.Where("debug_interface_id=?", debugInterfaceId)
	} else {
		db.Where("endpoint_interface_id=? AND debug_interface_id=?", endpointInterfaceId, 0)
	}

	if typ == consts.ConditionCategoryAssert {
		db.Where("entity_type = ?", consts.ConditionTypeCheckpoint)
	} else if typ == consts.ConditionCategoryConsole {
		db.Where("entity_type =? or entity_type =? or entity_type=?", consts.ConditionTypeExtractor, consts.ConditionTypeCookie, consts.ConditionTypeScript)
	} else if typ == consts.ConditionCategoryResponse {
		db.Where("entity_type = ?", consts.ConditionTypeResponseDefine)
	} else if typ == consts.ConditionCategoryResult {
		db.Where("entity_type = ? or entity_type = ?", consts.ConditionTypeResponseDefine, consts.ConditionTypeCheckpoint)
	}

	err = db.Find(&pos).Error

	return
}

func (r *PostConditionRepo) ListExtractor(debugInterfaceId, endpointInterfaceId uint) (pos []model.DebugPostCondition, err error) {
	db := r.DB.
		Where("NOT deleted").
		Order("ordr ASC")

	if debugInterfaceId > 0 {
		db.Where("debug_interface_id=?", debugInterfaceId)
	} else {
		db.Where("endpoint_interface_id=? AND debug_interface_id=?", endpointInterfaceId, 0)
	}

	db.Where("entity_type = ?", consts.ConditionTypeExtractor)

	err = db.Find(&pos).Error

	return
}

func (r *PostConditionRepo) Get(id uint) (po model.DebugPostCondition, err error) {
	err = r.DB.
		Where("id=?", id).
		Where("NOT deleted").
		First(&po).Error
	return
}

func (r *PostConditionRepo) Save(po *model.DebugPostCondition) (err error) {
	err = r.DB.Save(po).Error
	return
}

func (r *PostConditionRepo) CloneAll(srcDebugInterfaceId, srcEndpointInterfaceId, distDebugInterfaceId uint) (err error) {
	srcConditions, err := r.List(srcDebugInterfaceId, srcEndpointInterfaceId, consts.ConditionCategoryAll)

	for _, srcCondition := range srcConditions {
		// clone condition po
		srcCondition.ID = 0
		srcCondition.DebugInterfaceId = distDebugInterfaceId

		r.Save(&srcCondition)

		// clone condition entity
		var entityId uint
		if srcCondition.EntityType == consts.ConditionTypeExtractor {
			srcEntity, _ := r.ExtractorRepo.Get(srcCondition.EntityId)
			srcEntity.ID = 0
			srcEntity.ConditionId = srcCondition.ID

			r.ExtractorRepo.Save(&srcEntity)
			entityId = srcEntity.ID

		} else if srcCondition.EntityType == consts.ConditionTypeCookie {
			srcEntity, _ := r.CookieRepo.Get(srcCondition.EntityId)
			srcEntity.ID = 0
			srcEntity.ConditionId = srcCondition.ID

			r.CookieRepo.Save(&srcEntity)
			entityId = srcEntity.ID

		} else if srcCondition.EntityType == consts.ConditionTypeCheckpoint {
			srcEntity, _ := r.CheckpointRepo.Get(srcCondition.EntityId)
			srcEntity.ID = 0
			srcEntity.ConditionId = srcCondition.ID

			r.CheckpointRepo.Save(&srcEntity)
			entityId = srcEntity.ID

		} else if srcCondition.EntityType == consts.ConditionTypeScript {
			srcEntity, _ := r.ScriptRepo.Get(srcCondition.EntityId)
			srcEntity.ID = 0
			srcEntity.ConditionId = srcCondition.ID

			r.ScriptRepo.Save(&srcEntity)
			entityId = srcEntity.ID
		}

		err = r.UpdateEntityId(srcCondition.ID, entityId)
	}

	return
}

func (r *PostConditionRepo) ReplaceAll(debugInterfaceId, endpointInterfaceId uint, postConditions []domain.InterfaceExecCondition) (err error) {
	r.removeAll(debugInterfaceId, endpointInterfaceId)

	for _, item := range postConditions {
		// clone condition po
		condition := model.DebugPostCondition{
			EntityType:          item.Type,
			DebugInterfaceId:    debugInterfaceId,
			EndpointInterfaceId: endpointInterfaceId,
			Desc:                item.Desc,
		}
		r.Save(&condition)

		// clone condition entity
		var entityId uint
		if item.Type == consts.ConditionTypeExtractor {
			extractor := domain.ExtractorBase{}
			json.Unmarshal(item.Raw, &extractor)

			entity := model.DebugConditionExtractor{}

			copier.CopyWithOption(&entity, extractor, copier.Option{DeepCopy: true})
			entity.ID = 0
			entity.ConditionId = condition.ID

			r.ExtractorRepo.Save(&entity)
			entityId = entity.ID

		} else if item.Type == consts.ConditionTypeCookie {
			cookie := domain.CookieBase{}
			json.Unmarshal(item.Raw, &cookie)

			entity := model.DebugConditionCookie{}

			copier.CopyWithOption(&entity, cookie, copier.Option{DeepCopy: true})
			entity.ID = 0
			entity.ConditionId = condition.ID

			r.CookieRepo.Save(&entity)
			entityId = entity.ID

		} else if item.Type == consts.ConditionTypeCheckpoint {
			checkpoint := domain.CheckpointBase{}
			json.Unmarshal(item.Raw, &checkpoint)

			entity := model.DebugConditionCheckpoint{}

			copier.CopyWithOption(&entity, checkpoint, copier.Option{DeepCopy: true})
			entity.ID = 0
			entity.ConditionId = condition.ID

			r.CheckpointRepo.Save(&entity)
			entityId = entity.ID

		} else if item.Type == consts.ConditionTypeScript {
			script := domain.ScriptBase{}
			json.Unmarshal(item.Raw, &script)

			entity := model.DebugConditionScript{}

			copier.CopyWithOption(&entity, script, copier.Option{DeepCopy: true})
			entity.ID = 0
			entity.ConditionId = condition.ID

			r.ScriptRepo.Save(&entity)
			entityId = entity.ID
		}

		err = r.UpdateEntityId(condition.ID, entityId)
	}

	return
}

func (r *PostConditionRepo) Delete(id uint) (err error) {
	po, _ := r.Get(id)

	err = r.DB.Model(&model.DebugPostCondition{}).
		Where("id=?", id).
		Update("deleted", true).
		Error

	if po.EntityType == consts.ConditionTypeExtractor {
		r.ExtractorRepo.DeleteByCondition(id)
	} else if po.EntityType == consts.ConditionTypeCookie {
		r.CookieRepo.DeleteByCondition(id)
	} else if po.EntityType == consts.ConditionTypeCheckpoint {
		r.CheckpointRepo.DeleteByCondition(id)
	} else if po.EntityType == consts.ConditionTypeScript {
		r.ScriptRepo.DeleteByCondition(id)
	}

	return
}

func (r *PostConditionRepo) Disable(id uint) (err error) {
	err = r.DB.Model(&model.DebugPostCondition{}).
		Where("id=?", id).
		Update("disabled", gorm.Expr("NOT disabled")).
		Error

	return
}

func (r *PostConditionRepo) UpdateOrders(req serverDomain.ConditionMoveReq) (err error) {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		for index, id := range req.Data {
			sql := fmt.Sprintf("UPDATE %s SET ordr = %d WHERE id = %d",
				model.DebugPostCondition{}.TableName(), index+1, id)

			err = r.DB.Exec(sql).Error
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *PostConditionRepo) UpdateEntityId(id uint, entityId uint) (err error) {
	err = r.DB.Model(&model.DebugPostCondition{}).
		Where("id=?", id).
		Update("entity_id", entityId).
		Error

	return
}

func (r *PostConditionRepo) ListTo(debugInterfaceId, endpointInterfaceId uint) (ret []domain.InterfaceExecCondition, err error) {
	pos, err := r.List(debugInterfaceId, endpointInterfaceId, consts.ConditionCategoryAll)

	for _, po := range pos {
		typ := po.EntityType

		if typ == consts.ConditionTypeExtractor {
			extractor := domain.ExtractorBase{}

			entity, _ := r.ExtractorRepo.Get(po.EntityId)
			copier.CopyWithOption(&extractor, entity, copier.Option{DeepCopy: true})
			extractor.ConditionEntityType = typ
			extractor.ConditionId = po.ID
			extractor.ConditionEntityId = po.EntityId

			raw, _ := json.Marshal(extractor)
			condition := domain.InterfaceExecCondition{
				Type: typ,
				Raw:  raw,
				Desc: po.Desc,
			}

			ret = append(ret, condition)

		} else if typ == consts.ConditionTypeCookie {
			cookie := domain.CookieBase{}

			entity, _ := r.CookieRepo.Get(po.EntityId)
			copier.CopyWithOption(&cookie, entity, copier.Option{DeepCopy: true})
			cookie.ConditionEntityType = typ
			cookie.ConditionId = po.ID
			cookie.ConditionEntityId = po.EntityId

			raw, _ := json.Marshal(cookie)
			condition := domain.InterfaceExecCondition{
				Type: typ,
				Raw:  raw,
				Desc: po.Desc,
			}

			ret = append(ret, condition)

		} else if typ == consts.ConditionTypeCheckpoint {
			checkpoint := domain.CheckpointBase{}

			entity, _ := r.CheckpointRepo.Get(po.EntityId)
			copier.CopyWithOption(&checkpoint, entity, copier.Option{DeepCopy: true})
			checkpoint.ConditionEntityType = typ
			checkpoint.ConditionId = po.ID
			checkpoint.ConditionEntityId = po.EntityId

			raw, _ := json.Marshal(checkpoint)
			condition := domain.InterfaceExecCondition{
				Type: typ,
				Raw:  raw,
				Desc: po.Desc,
			}

			ret = append(ret, condition)

		} else if typ == consts.ConditionTypeScript {
			script := domain.ScriptBase{}

			entity, _ := r.ScriptRepo.Get(po.EntityId)
			copier.CopyWithOption(&script, entity, copier.Option{DeepCopy: true})
			script.ConditionId = po.ID
			script.ConditionEntityId = po.EntityId

			raw, _ := json.Marshal(script)
			condition := domain.InterfaceExecCondition{
				Type: typ,
				Raw:  raw,
			}

			ret = append(ret, condition)

		} else if typ == consts.ConditionTypeResponseDefine {
			responseDefine := domain.ResponseDefineBase{}

			entity, _ := r.ResponseDefineRepo.Get(po.EntityId)
			copier.CopyWithOption(&responseDefine, entity, copier.Option{DeepCopy: true})
			responseDefine.ConditionId = po.ID
			responseDefine.ConditionEntityId = po.EntityId
			responseBody := r.EndpointInterfaceRepo.GetResponse(endpointInterfaceId, entity.Code)
			responseDefine.Schema = responseBody.SchemaItem.Content
			responseDefine.Code = entity.Code
			responseDefine.MediaType = responseBody.MediaType
			components := r.ResponseDefineRepo.Components(endpointInterfaceId)
			responseDefine.Component = commonUtils.JsonEncode(components)
			raw, _ := json.Marshal(responseDefine)
			condition := domain.InterfaceExecCondition{
				Type: typ,
				Raw:  raw,
			}

			ret = append(ret, condition)
		}

	}

	return
}

func (r *PostConditionRepo) removeAll(debugInterfaceId, endpointInterfaceId uint) (err error) {
	pos, _ := r.List(debugInterfaceId, endpointInterfaceId, "")

	for _, po := range pos {
		r.Delete(po.ID)
	}

	return
}

func (r *PostConditionRepo) CreateDefaultResponseDefine(debugInterfaceId, endpointInterfaceId uint, by consts.UsedBy) (condition domain.Condition) {

	if endpointInterfaceId == 0 {
		return
	}

	po, err := r.GetByDebugInterfaceId(debugInterfaceId, endpointInterfaceId, by)
	if err == gorm.ErrRecordNotFound {
		po, err = r.saveDefault(debugInterfaceId, endpointInterfaceId, by)
		if err != nil {
			return
		}
	}

	copier.CopyWithOption(&condition, po, copier.Option{DeepCopy: true})

	entityData, _ := r.ResponseDefineRepo.Get(po.EntityId)
	entityData.Codes = r.EndpointInterfaceRepo.GetResponseCodes(endpointInterfaceId)
	//entityData.Component = r.ResponseDefineRepo.Components(endpointInterfaceId)
	condition.EntityData = entityData

	return
}

func (r *PostConditionRepo) GetByDebugInterfaceId(debugInterfaceId, endpointInterfaceId uint, by consts.UsedBy) (po model.DebugPostCondition, err error) {
	err = r.DB.
		Where("debug_interface_id=? and endpoint_interface_id=? and used_by=? and entity_type=?", debugInterfaceId, endpointInterfaceId, by, consts.ConditionTypeResponseDefine).
		Where("NOT deleted").
		First(&po).Error
	return
}

func (r *PostConditionRepo) saveDefault(debugInterfaceId, endpointInterfaceId uint, by consts.UsedBy) (po model.DebugPostCondition, err error) {

	responseDefine := model.DebugConditionResponseDefine{}
	responseDefine.Code = "200"
	err = r.ResponseDefineRepo.Save(&responseDefine)
	if err != nil {
		return
	}

	po.EntityType = consts.ConditionTypeResponseDefine
	po.EndpointInterfaceId = endpointInterfaceId
	po.DebugInterfaceId = debugInterfaceId
	po.UsedBy = by
	po.EntityId = responseDefine.ID
	err = r.Save(&po)

	return
}
