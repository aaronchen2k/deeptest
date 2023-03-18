package router

import (
	"github.com/aaronchen2k/deeptest/cmd/server/v1/handler"
	"github.com/aaronchen2k/deeptest/internal/pkg/core/module"
	"github.com/aaronchen2k/deeptest/internal/server/middleware"
	"github.com/kataras/iris/v12"
)

type CategoryModule struct {
	CategoryCtrl *handler.CategoryCtrl `inject:""`
}

// Party 场景
func (m *CategoryModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck(), middleware.JwtHandler(), middleware.OperationRecord(), middleware.Casbin())

		index.Get("/load", m.CategoryCtrl.LoadTree).Name = "分类树状数据"
		index.Get("/{id:uint}", m.CategoryCtrl.Get).Name = "分类详情"
		index.Post("/", m.CategoryCtrl.Create).Name = "新建分类"
		index.Put("/", m.CategoryCtrl.Update).Name = "更新分类"
		index.Put("/{id:uint}/updateName", m.CategoryCtrl.UpdateName).Name = "更新节点名称"
		index.Delete("/{id:uint}", m.CategoryCtrl.Delete).Name = "删除节点"
		index.Post("/move", m.CategoryCtrl.Move).Name = "移动节点"
	}

	return module.NewModule("/categories", handler)
}
