package router

import (
	"github.com/aaronchen2k/deeptest/cmd/server/v1/handler"
	"github.com/aaronchen2k/deeptest/internal/pkg/core/module"
	"github.com/aaronchen2k/deeptest/internal/server/middleware"
	"github.com/kataras/iris/v12"
)

type ExtractorModule struct {
	ExtractorCtrl *handler.ExtractorCtrl `inject:""`
}

// Party 提取器
func (m *ExtractorModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck(), middleware.JwtHandler(), middleware.OperationRecord(), middleware.Casbin())

		index.Get("/", m.ExtractorCtrl.List).Name = "提取器列表"
		index.Get("/{id:uint}", m.ExtractorCtrl.Get).Name = "提取器详情"
		//index.Post("/createOrUpdateResult", m.ExtractorCtrl.CreateOrUpdateResult).Name = "新建提取器"
		index.Put("/", m.ExtractorCtrl.Update).Name = "更新提取器"
		index.Delete("/{id:uint}", m.ExtractorCtrl.Delete).Name = "删除提取器"

		index.Post("/listExtractorVariableForCheckpoint", m.ExtractorCtrl.ListExtractorVariableForCheckpoint).Name = "提取器变量列表"
	}

	return module.NewModule("/extractors", handler)
}
