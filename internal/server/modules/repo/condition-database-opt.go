package repo

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	databaseOptHelpper "github.com/aaronchen2k/deeptest/internal/pkg/helper/database-opt"
	model "github.com/aaronchen2k/deeptest/internal/server/modules/model"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type DatabaseOptRepo struct {
	DB *gorm.DB `inject:""`
}

func (r *DatabaseOptRepo) Get(id uint) (databaseOpt model.DebugConditionDatabaseOpt, err error) {
	err = r.DB.
		Where("id=?", id).
		Where("NOT deleted").
		First(&databaseOpt).Error
	return
}

func (r *DatabaseOptRepo) Save(databaseOpt *model.DebugConditionDatabaseOpt) (err error) {
	r.UpdateDesc(databaseOpt)

	err = r.DB.Save(databaseOpt).Error
	return
}
func (r *DatabaseOptRepo) UpdateDesc(po *model.DebugConditionDatabaseOpt) (err error) {
	desc := databaseOptHelpper.GenDesc(po.Type, po.Sql)
	values := map[string]interface{}{
		"desc": desc,
	}

	err = r.DB.Model(&model.DebugPostCondition{}).
		Where("id=?", po.ConditionId).
		Updates(values).Error

	return
}

func (r *DatabaseOptRepo) Delete(id uint) (err error) {
	err = r.DB.Model(&model.DebugConditionDatabaseOpt{}).
		Where("id=?", id).
		Update("deleted", true).
		Error

	return
}
func (r *DatabaseOptRepo) DeleteByCondition(conditionId uint) (err error) {
	err = r.DB.Model(&model.DebugConditionDatabaseOpt{}).
		Where("condition_id=?", conditionId).
		Update("deleted", true).
		Error

	return
}

func (r *DatabaseOptRepo) UpdateResult(databaseOpt domain.DatabaseOptBase) (err error) {
	values := map[string]interface{}{
		"result_msg":    databaseOpt.ResultMsg,
		"result_status": databaseOpt.ResultStatus,
	}

	err = r.DB.Model(&model.DebugConditionDatabaseOpt{}).
		Where("id=?", databaseOpt.ConditionEntityId).
		Updates(values).
		Error

	return
}

func (r *DatabaseOptRepo) CreateLog(databaseOpt domain.DatabaseOptBase) (
	log model.ExecLogDatabaseOpt, err error) {

	copier.CopyWithOption(&log, databaseOpt, copier.Option{DeepCopy: true})

	log.ID = 0
	log.ConditionId = databaseOpt.ConditionId
	log.ConditionEntityId = databaseOpt.ConditionEntityId

	log.InvokeId = databaseOpt.InvokeId
	log.CreatedAt = nil
	log.UpdatedAt = nil

	err = r.DB.Save(&log).Error

	return
}

func (r *DatabaseOptRepo) CreateDefault(conditionId uint) (po model.DebugConditionDatabaseOpt) {
	po.ConditionId = conditionId

	po = model.DebugConditionDatabaseOpt{
		DatabaseOptBase: domain.DatabaseOptBase{
			ConditionId: conditionId,

			DatabaseConnBase: domain.DatabaseConnBase{
				Type: consts.DbTypeMySql,
			},
		},
	}

	r.Save(&po)

	return
}

func (r *DatabaseOptRepo) GetLog(conditionId, invokeId uint) (ret model.ExecLogDatabaseOpt, err error) {
	err = r.DB.
		Where("condition_id=? AND invoke_id=?", conditionId, invokeId).
		Where("NOT deleted").
		First(&ret).Error

	ret.ConditionEntityType = consts.ConditionTypeDatabase

	return
}