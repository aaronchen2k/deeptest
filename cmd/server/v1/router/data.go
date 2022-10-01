package router

import (
	"github.com/aaronchen2k/deeptest/cmd/server/v1/handler"
	"github.com/aaronchen2k/deeptest/internal/server/core/module"
	"github.com/kataras/iris/v12"
)

type DataModule struct {
	DataCtrl *handler.DataCtrl `inject:""`
}

func NewDataModule() *DataModule {
	return &DataModule{}
}

// Party 初始化模块
func (m *DataModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Post("/initdb", m.DataCtrl.Init)
		index.Get("/checkdb", m.DataCtrl.Check)
	}
	return module.NewModule("/init", handler)
}
