package repo

import (
	v1 "github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	serverConsts "github.com/aaronchen2k/deeptest/internal/server/consts"
	"github.com/aaronchen2k/deeptest/internal/server/modules/model"
	"github.com/kataras/iris/v12"
	"gorm.io/gorm"
)

type CategoryRepo struct {
	DB          *gorm.DB     `inject:""`
	ProjectRepo *ProjectRepo `inject:""`
}

func (r *CategoryRepo) GetTree(typ serverConsts.CategoryDiscriminator, projectId, serveId uint) (root *v1.Category, err error) {
	pos, err := r.ListByProject(typ, projectId, serveId)
	if err != nil {
		return
	}

	tos := r.toTos(pos)
	if len(tos) == 0 {
		return
	}

	root = tos[0]
	root.Slots = iris.Map{"icon": "icon"}

	r.makeTree(tos[1:], root)

	return
}

func (r *CategoryRepo) ListByProject(typ serverConsts.CategoryDiscriminator, projectId, serveId uint) (pos []*model.Category, err error) {
	db := r.DB.
		Where("project_id=?", projectId).
		Where("type=?", typ).
		Where("NOT deleted")

	if serveId > 0 {
		db.Where("serve_id=?", serveId)
	}

	err = db.
		Order("parent_id ASC, ordr ASC").
		Find(&pos).Error

	return
}

func (r *CategoryRepo) Get(id uint) (po model.Category, err error) {
	err = r.DB.Where("id = ?", id).First(&po).Error
	return
}

func (r *CategoryRepo) toTos(pos []*model.Category) (tos []*v1.Category) {
	for _, po := range pos {
		to := v1.Category{
			Id:       int64(po.ID),
			Name:     po.Name,
			Desc:     po.Desc,
			ParentId: int64(po.ParentId),
		}

		tos = append(tos, &to)
	}

	return
}

func (r *CategoryRepo) makeTree(findIn []*v1.Category, parent *v1.Category) { //参数为父节点，添加父节点的子节点指针切片
	children, _ := r.hasChild(findIn, parent) // 判断节点是否有子节点并返回

	if children != nil {
		parent.Children = append(parent.Children, children[0:]...) // 添加子节点

		for _, child := range children { // 查询子节点的子节点，并添加到子节点
			_, has := r.hasChild(findIn, child)
			if has {
				r.makeTree(findIn, child) // 递归添加节点
			}
		}
	}
}

func (r *CategoryRepo) hasChild(categories []*v1.Category, parent *v1.Category) (
	ret []*v1.Category, yes bool) {

	for _, item := range categories {
		if item.ParentId == parent.Id {
			item.Slots = iris.Map{"icon": "icon"}
			//item.Parent = parent // loop json

			ret = append(ret, item)
		}
	}

	if ret != nil {
		yes = true
	}

	return
}

func (r *CategoryRepo) Save(processor *model.Category) (err error) {
	err = r.DB.Save(processor).Error

	return
}

func (r *CategoryRepo) UpdateOrder(pos serverConsts.DropPos, targetId uint, typ serverConsts.CategoryDiscriminator, projectId uint) (
	parentId uint, ordr int) {
	if pos == serverConsts.Inner {
		parentId = targetId

		var preChild model.Category
		r.DB.Where("parent_id=? AND type = ? AND project_id = ?", parentId, typ, projectId).
			Order("ordr DESC").Limit(1).
			First(&preChild)

		ordr = preChild.Ordr + 1

	} else if pos == serverConsts.Before {
		brother, _ := r.Get(targetId)
		parentId = brother.ParentId

		r.DB.Model(&model.Category{}).
			Where("NOT deleted AND parent_id=? AND type = ? AND project_id = ? AND ordr >= ?",
				parentId, typ, projectId, brother.Ordr).
			Update("ordr", gorm.Expr("ordr + 1"))

		ordr = brother.Ordr

	} else if pos == serverConsts.After {
		brother, _ := r.Get(targetId)
		parentId = brother.ParentId

		r.DB.Model(&model.Category{}).
			Where("NOT deleted AND parent_id=? AND type = ? AND project_id = ? AND ordr > ?",
				parentId, parentId, typ, brother.Ordr).
			Update("ordr", gorm.Expr("ordr + 1"))

		ordr = brother.Ordr + 1

	}

	return
}

func (r *CategoryRepo) UpdateName(id int, name string) (err error) {
	err = r.DB.Model(&model.Category{}).
		Where("id = ?", id).
		Update("name", name).Error

	return
}

func (r *CategoryRepo) Update(req v1.CategoryReq) (err error) {
	po := new(model.Category)
	po.ID = uint(req.Id)

	err = r.DB.First(&po).Error
	if err != nil {
		return err
	}

	po.Name = req.Name
	po.Desc = req.Desc

	err = r.DB.Save(&po).Error

	return
}

func (r *CategoryRepo) Delete(id uint) (err error) {
	err = r.DB.Model(&model.Category{}).
		Where("id=?", id).
		Update("deleted", true).
		Error

	return
}

func (r *CategoryRepo) GetChildren(nodeId uint) (children []*model.Category, err error) {
	err = r.DB.Where("parent_id=?", nodeId).Find(&children).Error
	return
}

func (r *CategoryRepo) UpdateOrdAndParent(node model.Category) (err error) {
	err = r.DB.Model(&node).
		Updates(model.Category{Ordr: node.Ordr, ParentId: node.ParentId}).
		Error

	return
}

func (r *CategoryRepo) GetMaxOrder(parentId uint, typ serverConsts.CategoryDiscriminator, projectId uint) (order int) {
	node := model.Category{}

	err := r.DB.Model(&model.Category{}).
		Where("parent_id=? AND type = ? AND project_id = ?", parentId, typ, projectId).
		Order("ordr DESC").
		First(&node).Error

	if err == nil {
		order = node.Ordr + 1
	}

	return
}
