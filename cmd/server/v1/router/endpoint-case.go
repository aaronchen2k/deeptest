package router

import (
	"github.com/aaronchen2k/deeptest/cmd/server/v1/handler"
	"github.com/aaronchen2k/deeptest/internal/pkg/core/module"
	"github.com/aaronchen2k/deeptest/internal/server/middleware"
	"github.com/kataras/iris/v12"
)

type EndpointCaseModule struct {
	EndpointCaseCtrl *handler.EndpointCaseCtrl `inject:""`
}

// Party 注册模块
func (m *EndpointCaseModule) Party() module.WebModule {
	handler := func(public iris.Party) {
		public.Use(middleware.InitCheck(), middleware.JwtHandler(), middleware.OperationRecord(), middleware.Casbin())

		public.Get("/list", m.EndpointCaseCtrl.List).Name = "用例列表"
		public.Get("/{id:uint}", m.EndpointCaseCtrl.Get).Name = "用例详情"
		public.Post("/{id:uint}", m.EndpointCaseCtrl.Save).Name = "保存用例"
		public.Put("/updateName", m.EndpointCaseCtrl.UpdateName).Name = "保存用例名称"
		public.Post("/saveDebugData", m.EndpointCaseCtrl.SaveDebugData).Name = "保存调试数据"
		public.Delete("/{id:uint}", m.EndpointCaseCtrl.Remove).Name = "删除用例"
	}

	return module.NewModule("/endpoints/cases", handler)
}
