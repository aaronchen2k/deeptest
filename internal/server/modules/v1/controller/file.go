package controller

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
	"github.com/kataras/iris/v12"
	"go.uber.org/zap"
)

type FileCtrl struct {
	FileService *service.FileService `inject:""`
}

func NewFileCtrl() *FileCtrl {
	return &FileCtrl{}
}

// Upload 上传文件
// - 需要 file 表单文件字段
func (c *FileCtrl) Upload(ctx iris.Context) {
	f, fh, err := ctx.FormFile("file")
	if err != nil {
		logUtils.Errorf("文件上传失败", zap.String("ctx.FormFile(\"file\")", err.Error()))
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}
	defer f.Close()

	pth, err := c.FileService.UploadFile(ctx, fh)
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}
	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: pth, Msg: domain.NoErr.Msg})
}
