package repo

import (
	"errors"
	"fmt"
	v1 "github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	casbinServer "github.com/aaronchen2k/deeptest/internal/server/core/casbin"
	"github.com/aaronchen2k/deeptest/internal/server/core/dao"
	"github.com/aaronchen2k/deeptest/internal/server/modules/model"
	"github.com/aaronchen2k/deeptest/pkg/domain"
	logUtils "github.com/aaronchen2k/deeptest/pkg/lib/log"
	"strconv"

	"gorm.io/gorm"
)

type PermRepo struct {
	DB *gorm.DB `inject:""`

	RoleRepo *RoleRepo `inject:""`
}

func NewPermRepo() *PermRepo {
	return &PermRepo{}
}

// Paginate
func (r *PermRepo) Paginate(req v1.PermReqPaginate) (data _domain.PageData, err error) {
	var count int64
	db := r.DB.Model(&model.SysPerm{})
	if req.Name != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%s%%", req.Name))
	}

	err = db.Count(&count).Error
	if err != nil {
		logUtils.Errorf("获取权限总数失败, 错误%s。", err.Error())
		return
	}

	perms := make([]*v1.PermResp, 0)
	err = db.Scopes(dao.PaginateScope(req.Page, req.PageSize, req.Order, req.Field)).Find(&perms).Error
	if err != nil {
		logUtils.Errorf("获取权限分页数据失败, 错误%s。", err.Error())
		return
	}

	data.Populate(perms, count, req.Page, req.PageSize)

	return
}

// FindByNameAndAct
// db *gorm.DB
// name 名称
// act 方法
// ids 当 ids 的 len = 1 ，排除次 id 数据
func (r *PermRepo) FindByNameAndAct(name, act string, ids ...uint) (v1.PermResp, error) {
	perm := v1.PermResp{}
	db := r.DB.Model(&model.SysPerm{}).Where("name = ?", name).Where("act = ?", act)
	if len(ids) == 1 {
		db.Where("id != ?", ids[0])
	}
	err := db.First(&perm).Error
	if err != nil {
		logUtils.Errorf("根据名称和方法获取权限数据失败, 错误%s。", err.Error())
		return perm, err
	}
	return perm, nil
}

// Create
func (r *PermRepo) Create(req v1.PermReq) (uint, error) {
	perm := model.SysPerm{PermBase: req.PermBase}
	if !r.CheckNameAndAct(req) {
		return perm.ID, fmt.Errorf("权限[%s-%s]已存在", req.Name, req.Act)
	}
	err := r.DB.Model(&model.SysPerm{}).Create(&perm).Error
	if err != nil {
		logUtils.Errorf("添加权限失败，错误%s。", err.Error())
		return perm.ID, err
	}
	return perm.ID, nil
}

// CreateInBatches
func (r *PermRepo) CreateInBatches(perms []model.SysPerm) error {
	err := r.DB.Model(&model.SysPerm{}).CreateInBatches(&perms, 500).Error
	if err != nil {
		logUtils.Errorf("添加权限失败，错误%s。", err.Error())
		return err
	}
	return nil
}

// CreateIfNotExist
func (r *PermRepo) CreateIfNotExist(perms []model.SysPerm) (count int, err error) {
	enforcer := casbinServer.Instance()

	adminRole, _ := r.RoleRepo.FindFirstAdminUser()
	adminRoleId := strconv.Itoa(int(adminRole.Id))

	r.DB.Transaction(func(tx *gorm.DB) (err error) {
		for _, perm := range perms {
			found := enforcer.HasNamedPolicy("p", adminRoleId, perm.Name, perm.Act)
			if found {
				continue
			}

			// add to casbin table
			namedPolicy := []string{adminRoleId, perm.Name, perm.Act}
			success, _ := enforcer.AddNamedPolicy("p", namedPolicy)
			if success {
				count++
			}

			// add to permission table
			err = r.DB.Model(&model.SysPerm{}).Create(&perm).Error
			if err != nil {
				logUtils.Errorf("添加权限%s失败，错误%s。", perm.Name, err.Error())
				continue
			}
		}

		return
	})

	return
}

// Update
func (r *PermRepo) Update(id uint, req v1.PermReq) error {
	if !r.CheckNameAndAct(req, id) {
		return fmt.Errorf("权限[%s-%s]已存在", req.Name, req.Act)
	}
	perm := model.SysPerm{PermBase: req.PermBase}
	err := r.DB.Model(&model.SysPerm{}).Where("id = ?", id).Updates(&perm).Error
	if err != nil {
		logUtils.Errorf("更新权限失败, 错误%s。", err.Error())
		return err
	}
	return nil
}

// checkNameAndAct
func (r *PermRepo) CheckNameAndAct(req v1.PermReq, ids ...uint) bool {
	_, err := r.FindByNameAndAct(req.Name, req.Act, ids...)
	return errors.Is(err, gorm.ErrRecordNotFound)
}

// FindById
func (r *PermRepo) FindById(id uint) (v1.PermResp, error) {
	res := v1.PermResp{}
	err := r.DB.Model(&model.SysPerm{}).Where("id = ?", id).First(&res).Error
	if err != nil {
		logUtils.Errorf("获取权限失败, 错误%s。", err.Error())
		return res, err
	}
	return res, nil
}

// DeleteById
func (r *PermRepo) DeleteById(id uint) error {
	err := r.DB.Unscoped().Delete(&model.SysPerm{}, id).Error
	if err != nil {
		logUtils.Errorf("删除权限失败, 错误%s。", err.Error())
		return err
	}
	return nil
}

// DeleteAll, for init
func (r *PermRepo) DeleteAll() error {
	err := r.DB.Where("1 = 1").Delete(&model.SysPerm{}).Error
	if err != nil {
		logUtils.Errorf("删除权限失败, 错误%s。", err.Error())
		return err
	}
	return nil
}

// GetPermsForRole
func (r *PermRepo) GetPermsForRole() ([][]string, error) {
	var permsForRoles [][]string
	var perms []model.SysPerm
	err := r.DB.Model(&model.SysPerm{}).Find(&perms).Error
	if err != nil {
		return nil, fmt.Errorf("获取权限错误 %w", err)
	}
	for _, perm := range perms {
		permsForRole := []string{perm.Name, perm.Act}
		permsForRoles = append(permsForRoles, permsForRole)
	}
	return permsForRoles, nil
}

func (r *PermRepo) GetPermsForRoles() (map[consts.RoleType][][]string, error) {
	permsUserNotInclude := []v1.PermStruct{
		{
			Name: "/api/v1/users",
			Act:  "POST",
		},
		{
			Name: "/api/v1/users/:id",
			Act:  "POST",
		},
		{
			Name: "/api/v1/users/:id",
			Act:  "DELETE",
		},
	}
	var adminPerms, userPerms [][]string
	var perms []model.SysPerm
	err := r.DB.Model(&model.SysPerm{}).Find(&perms).Error
	if err != nil {
		return nil, fmt.Errorf("获取权限错误 %w", err)
	}

OuterLoop:
	for _, perm := range perms {
		permForRole := []string{perm.Name, perm.Act}
		adminPerms = append(adminPerms, permForRole)
		for _, v := range permsUserNotInclude {
			if perm.Name == v.Name && perm.Act == v.Act {
				continue OuterLoop
			}
		}
		userPerms = append(userPerms, permForRole)
	}

	rolePermsMap := make(map[consts.RoleType][][]string)
	rolePermsMap[consts.Admin] = adminPerms
	rolePermsMap[consts.User] = userPerms
	return rolePermsMap, nil
}
