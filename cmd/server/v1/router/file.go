package router

import (
	"github.com/aaronchen2k/deeptest/cmd/server/v1/handler"
	serverConfig "github.com/aaronchen2k/deeptest/internal/server/config"
	"github.com/aaronchen2k/deeptest/internal/server/core/module"
	"github.com/aaronchen2k/deeptest/internal/server/middleware"
	"github.com/kataras/iris/v12"
)

type FileModule struct {
	FileCtrl *handler.FileCtrl `inject:""`
}

func NewFileModule() *FileModule {
	return &FileModule{}
}

// Party 上传文件模块
func (m *FileModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck(), middleware.JwtHandler(), middleware.OperationRecord(), middleware.Casbin())
		index.Post("/", iris.LimitRequestBodySize(serverConfig.CONFIG.MaxSize+1<<20), m.FileCtrl.Upload).Name = "上传文件"
	}
	return module.NewModule("/upload", handler)
}
