package index

import (
	"github.com/aaronchen2k/deeptest/internal/server/core/module"
	"github.com/aaronchen2k/deeptest/internal/server/middleware"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/controller"
	"github.com/kataras/iris/v12"
)

type ProjectModule struct {
	ProjectCtrl *controller.ProjectCtrl `inject:""`
}

func NewProjectModule() *ProjectModule {
	return &ProjectModule{}
}

// Party 项目
func (m *ProjectModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck(), middleware.JwtHandler(), middleware.OperationRecord(), middleware.Casbin())
		index.Get("/", m.ProjectCtrl.List).Name = "项目列表"
		index.Get("/{id:uint}", m.ProjectCtrl.Get).Name = "项目详情"
		index.Post("/", m.ProjectCtrl.Create).Name = "创建项目"
		index.Post("/{id:uint}", m.ProjectCtrl.Update).Name = "编辑项目"
		index.Delete("/{id:uint}", m.ProjectCtrl.Delete).Name = "删除项目"
	}
	return module.NewModule("/projects", handler)
}
