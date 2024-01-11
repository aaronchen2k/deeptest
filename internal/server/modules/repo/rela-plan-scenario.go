package repo

import (
	"github.com/aaronchen2k/deeptest/internal/server/modules/model"
	"gorm.io/gorm"
)

type RelaPlanScenarioRepo struct {
	*BaseRepo `inject:""`
}

func (r *RelaPlanScenarioRepo) Get(id uint) (res model.RelaPlanScenario, err error) {
	err = r.DB.Where("id = ?", id).First(&res).Error
	return
}

func (r *RelaPlanScenarioRepo) UpdateOrdrById(id uint, ordr int) (err error) {
	err = r.DB.Model(model.RelaPlanScenario{}).Where("id=?", id).Update("ordr", ordr).Error
	return
}

func (r *RelaPlanScenarioRepo) IncreaseOrderAfter(id, planId uint) (err error) {
	err = r.DB.Model(model.RelaPlanScenario{}).Where("id >= ? and plan_id = ?  and not deleted", id, planId).UpdateColumn("ordr", gorm.Expr("ordr + ?", 1)).Error
	return
}

func (r *RelaPlanScenarioRepo) GetMaxOrder(planId uint) (order int) {
	res := model.RelaPlanScenario{}

	err := r.DB.Model(&model.RelaPlanScenario{}).
		Where("plan_id=? AND not deleted", planId).
		Order("ordr DESC").
		First(&res).Error

	if err == nil {
		order = res.Ordr + 1
	}

	return
}
