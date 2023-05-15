package repo

import (
	"fmt"
	v1 "github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	serverConsts "github.com/aaronchen2k/deeptest/internal/server/consts"
	"github.com/aaronchen2k/deeptest/internal/server/core/dao"
	"github.com/aaronchen2k/deeptest/internal/server/modules/model"
	"github.com/aaronchen2k/deeptest/pkg/domain"
	commonUtils "github.com/aaronchen2k/deeptest/pkg/lib/comm"
	logUtils "github.com/aaronchen2k/deeptest/pkg/lib/log"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
)

type ScenarioRepo struct {
	DB          *gorm.DB     `inject:""`
	BaseRepo    *BaseRepo    `inject:""`
	ProjectRepo *ProjectRepo `inject:""`
}

func NewScenarioRepo() *ScenarioRepo {
	return &ScenarioRepo{}
}

func (r *ScenarioRepo) ListByProject(projectId int) (pos []model.Scenario, err error) {
	err = r.DB.
		Where("project_id=?", projectId).
		Where("NOT deleted").
		Find(&pos).Error
	return
}

func (r *ScenarioRepo) Paginate(req v1.ScenarioReqPaginate, projectId int) (data _domain.PageData, err error) {
	var count int64
	var categoryIds []uint

	if req.CategoryId > 0 {
		categoryIds, err = r.BaseRepo.GetAllChildIds(uint(req.CategoryId), model.Category{}.TableName(),
			serverConsts.ScenarioCategory, projectId)
		if err != nil {
			return
		}
	}

	db := r.DB.Model(&model.Scenario{}).
		Where("project_id = ? AND NOT deleted",
			projectId)

	if len(categoryIds) > 0 {
		db.Where("category_id IN(?)", categoryIds)
	} else if req.CategoryId == -1 {
		db.Where("category_id IN(?)", -1)
	}

	if req.Keywords != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", req.Keywords))
	}
	if req.Enabled != "" {
		db = db.Where("disabled = ?", commonUtils.IsDisable(req.Enabled))
	}
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}
	if req.Priority != "" {
		db = db.Where("priority = ?", req.Priority)
	}
	if req.Type != "" {
		db = db.Where("type = ?", req.Type)
	}

	err = db.Count(&count).Error
	if err != nil {
		logUtils.Errorf("count scenario error", zap.String("error:", err.Error()))
		return
	}

	scenarios := make([]*model.Scenario, 0)
	req.Order = "desc"
	err = db.
		Scopes(dao.PaginateScope(req.Page, req.PageSize, req.Order, req.Field)).
		Find(&scenarios).Error
	if err != nil {
		logUtils.Errorf("query scenario error", zap.String("error:", err.Error()))
		return
	}

	data.Populate(scenarios, count, req.Page, req.PageSize)

	return
}

func (r *ScenarioRepo) Get(id uint) (scenario model.Scenario, err error) {
	err = r.DB.Model(&model.Scenario{}).Where("id = ?", id).First(&scenario).Error
	if err != nil {
		logUtils.Errorf("find scenario by id error", zap.String("error:", err.Error()))
		return scenario, err
	}

	return scenario, nil
}

func (r *ScenarioRepo) FindByName(scenarioName string, id uint) (scenario model.Scenario, err error) {
	db := r.DB.Model(&model.Scenario{}).
		Where("name = ? AND NOT deleted", scenarioName)

	if id > 0 {
		db.Where("id != ?", id)
	}

	db.First(&scenario)

	return
}

func (r *ScenarioRepo) Create(scenario model.Scenario) (ret model.Scenario, bizErr *_domain.BizErr) {
	//po, err := r.FindByName(scenario.Name, 0)
	//if po.Name != "" {
	//	bizErr = &_domain.BizErr{Code: _domain.ErrNameExist.Code}
	//	return
	//}

	err := r.DB.Model(&model.Scenario{}).Create(&scenario).Error
	if err != nil {
		logUtils.Errorf("add scenario error", zap.String("error:", err.Error()))
		bizErr = &_domain.BizErr{Code: _domain.SystemErr.Code}

		return
	}

	err = r.UpdateSerialNumber(scenario.ID, scenario.ProjectId)
	if err != nil {
		logUtils.Errorf("update scenario serial number error", zap.String("error:", err.Error()))
		bizErr = &_domain.BizErr{Code: _domain.SystemErr.Code}

		return
	}
	ret = scenario

	return
}

func (r *ScenarioRepo) Update(req model.Scenario) error {
	values := map[string]interface{}{
		"name":             req.Name,
		"desc":             req.Desc,
		"disabled":         req.Disabled,
		"create_user_id":   req.CreateUserId,
		"create_user_name": req.CreateUserName,
		"priority":         req.Priority,
		"type":             req.Type,
		"status":           req.Status,
	}
	err := r.DB.Model(&req).Where("id = ?", req.ID).Updates(values).Error
	if err != nil {
		logUtils.Errorf("update scenario error", zap.String("error:", err.Error()))
		return err
	}

	return nil
}

func (r *ScenarioRepo) DeleteById(id uint) (err error) {
	err = r.DB.Model(&model.Scenario{}).Where("id = ?", id).
		Updates(map[string]interface{}{"deleted": true}).Error
	if err != nil {
		logUtils.Errorf("delete scenario by id error", zap.String("error:", err.Error()))
		return
	}

	return
}

func (r *ScenarioRepo) DeleteChildren(ids []int, tx *gorm.DB) (err error) {
	err = tx.Model(&model.Scenario{}).Where("id IN (?)", ids).
		Updates(map[string]interface{}{"deleted": true}).Error
	if err != nil {
		logUtils.Errorf("batch delete scenario error", zap.String("error:", err.Error()))
		return err
	}

	return nil
}

func (r *ScenarioRepo) GetChildrenIds(id uint) (ids []int, err error) {
	tmpl := `
		WITH RECURSIVE scenario AS (
			SELECT * FROM biz_scenario WHERE id = %d
			UNION ALL
			SELECT child.* FROM biz_scenario child, scenario WHERE child.parent_id = scenario.id
		)
		SELECT id FROM scenario WHERE id != %d
    `
	sql := fmt.Sprintf(tmpl, id, id)
	err = r.DB.Raw(sql).Scan(&ids).Error
	if err != nil {
		logUtils.Errorf("get children scenario error", zap.String("error:", err.Error()))
		return
	}

	return
}

func (r *ScenarioRepo) UpdateSerialNumber(id, projectId uint) (err error) {
	var project model.Project
	project, err = r.ProjectRepo.Get(projectId)
	if err != nil {
		return
	}

	err = r.DB.Model(&model.Scenario{}).Where("id=?", id).Update("serial_number", project.ShortName+"-TS-"+strconv.Itoa(int(id))).Error
	return
}

func (r *ScenarioRepo) ListScenarioRelation(id uint) (pos []model.RelaPlanScenario, err error) {
	err = r.DB.
		Where("scenario_id=?", id).
		Where("NOT deleted").
		Find(&pos).Error
	return
}

func (r *ScenarioRepo) AddPlans(scenarioId uint, planIds []int) (err error) {
	relations, _ := r.ListScenarioRelation(scenarioId)
	existMap := map[uint]bool{}
	for _, item := range relations {
		existMap[item.ScenarioId] = true
	}

	var pos []model.RelaPlanScenario

	for _, id := range planIds {
		if existMap[uint(id)] {
			continue
		}

		po := model.RelaPlanScenario{
			PlanId:     uint(id),
			ScenarioId: scenarioId,
		}
		pos = append(pos, po)
	}

	err = r.DB.Create(&pos).Error

	return
}

func (r *ScenarioRepo) PlanList(req v1.ScenarioPlanReqPaginate, scenarioId int) (data _domain.PageData, err error) {
	relations, _ := r.ListScenarioRelation(uint(scenarioId))
	var planIds []uint
	for _, item := range relations {
		planIds = append(planIds, item.PlanId)
	}

	db := r.DB.Model(&model.Plan{}).Where("not deleted and project_id=?", req.ProjectId)

	if len(planIds) > 0 {
		if req.Ref {
			db = db.Where(" id in (?)", planIds)
		} else {
			db = db.Where(" id not in (?)", planIds)
		}
	}

	var count int64

	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}

	if req.UpdateUserId != 0 {
		db = db.Where("update_user_id = ?", req.UpdateUserId)
	}

	if req.Keywords != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", req.Keywords))
	}

	err = db.Count(&count).Error
	if err != nil {
		logUtils.Errorf("count plan error", zap.String("error:", err.Error()))
		return
	}

	plans := make([]*model.Plan, 0)

	err = db.
		Scopes(dao.PaginateScope(req.Page, req.PageSize, req.Order, req.Field)).
		Find(&plans).Error
	if err != nil {
		logUtils.Errorf("query plan error", zap.String("error:", err.Error()))
		return
	}

	data.Populate(plans, count, req.Page, req.PageSize)

	return
}

func (r *ScenarioRepo) UpdateStatus(id uint, status consts.TestStatus) error {
	return r.DB.Model(&model.Scenario{}).Where("id = ?", id).Update("status", status).Error
}

func (r *ScenarioRepo) UpdatePriority(id uint, priority string) error {
	return r.DB.Model(&model.Scenario{}).Where("id = ?", id).Update("priority", priority).Error
}

func (r *ScenarioRepo) GetByIds(ids []uint) (scenarios []model.Scenario, err error) {
	err = r.DB.Model(&model.Scenario{}).Where("id IN (?)", ids).Find(&scenarios).Error
	if err != nil {
		logUtils.Errorf("find scenario by id error", zap.String("error:", err.Error()))
		return scenarios, err
	}

	return
}

func (r *ScenarioRepo) RemovePlans(scenarioId uint, planIds []int) (err error) {
	err = r.DB.Model(&model.RelaPlanScenario{}).Where("scenario_id=? and plan_id in (?)", scenarioId, planIds).Update("deleted", true).Error
	return
}
