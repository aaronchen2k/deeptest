package index

import (
	"github.com/kataras/iris/v12"
	"github.com/aaronchen2k/deeptest/internal/server/core/module"
	"github.com/aaronchen2k/deeptest/internal/server/middleware"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/controller"
)

type ProductModule struct {
	ProductCtrl *controller.ProductCtrl `inject:""`
}

func NewProductModule() *ProductModule {
	return &ProductModule{}
}

// Party 产品
func (m *ProductModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck(), middleware.JwtHandler(), middleware.OperationRecord(), middleware.Casbin())
		index.Get("/", m.ProductCtrl.Query).Name = "产品查询"
		index.Get("/{id:uint}", m.ProductCtrl.Get).Name = "产品详情"
		index.Post("/", m.ProductCtrl.Create).Name = "创建产品"
		index.Post("/{id:uint}", m.ProductCtrl.Update).Name = "编辑产品"
		index.Delete("/{id:uint}", m.ProductCtrl.Delete).Name = "删除产品"
	}
	return module.NewModule("/products", handler)
}
