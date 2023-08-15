package router

import (
	"github.com/aaronchen2k/deeptest/cmd/server/v1/handler"
	"github.com/aaronchen2k/deeptest/internal/pkg/core/module"
	"github.com/aaronchen2k/deeptest/internal/server/middleware"
	"github.com/kataras/iris/v12"
)

type ScenarioReportModule struct {
	ScenarioReportCtrl *handler.ScenarioReportCtrl `inject:""`
}

// Party 场景
func (m *ScenarioReportModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck(), middleware.JwtHandler(), middleware.OperationRecord(), middleware.Casbin(), middleware.ProjectPerm())
		index.Post("/", m.ScenarioReportCtrl.List).Name = "结果列表"
		index.Get("/{id:uint}", m.ScenarioReportCtrl.Get).Name = "结果详情"
		index.Delete("/{id:uint}", m.ScenarioReportCtrl.Delete).Name = "删除报告"
		index.Put("/{id:uint}", m.ScenarioReportCtrl.Create).Name = "创建报告"
	}

	return module.NewModule("/scenarios/reports", handler)
}
