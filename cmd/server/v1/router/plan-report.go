package router

import (
	"github.com/aaronchen2k/deeptest/cmd/server/v1/handler"
	"github.com/aaronchen2k/deeptest/internal/pkg/core/module"
	"github.com/aaronchen2k/deeptest/internal/server/middleware"
	"github.com/kataras/iris/v12"
)

type PlanReportModule struct {
	ReportCtrl *handler.PlanReportCtrl `inject:""`
}

func NewReportModule() *PlanReportModule {
	return &PlanReportModule{}
}

// Party
func (m *PlanReportModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck(), middleware.JwtHandler(), middleware.OperationRecord(), middleware.Casbin())
		index.Get("/", m.ReportCtrl.List).Name = "结果列表"
		index.Get("/{id:uint}", m.ReportCtrl.Get).Name = "结果详情"
		index.Delete("/{id:uint}", m.ReportCtrl.Delete).Name = "删除场景"
	}

	return module.NewModule("/plans/reports", handler)
}
