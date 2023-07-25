package handler

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/service"
	"github.com/aaronchen2k/deeptest/internal/server/modules/service"
	"github.com/aaronchen2k/deeptest/pkg/domain"
	logUtils "github.com/aaronchen2k/deeptest/pkg/lib/log"
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/helper/dir"
	"go.uber.org/zap"
	"path/filepath"
)

type FileCtrl struct {
	FileService     *commService.FileService `inject:""`
	DatapoolService *service.DatapoolService `inject:""`
}

// Upload 上传文件
// @Tags	上传文件模块
// @summary	上传文件
// @Produce	application/json
// @Param 	Authorization	header		string	true	"Authentication header"
// @Param 	currProjectId	query		int		true	"当前项目ID"
// @Param 	isDatapool 		query 		bool 	true 	"是否是数据池"
// @Param 	file 			formData 	string true 	"文件"
// @success	200	{object}	_domain.Response{data=object{path=string,data=interface{}}}
// @Router	/api/v1/upload	[post]
func (c *FileCtrl) Upload(ctx iris.Context) {
	isDatapool, _ := ctx.URLParamBool("isDatapool")

	f, fh, err := ctx.FormFile("file")
	if err != nil {
		logUtils.Errorf("文件上传失败", zap.String("ctx.FormFile(\"file\")", err.Error()))
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}
	defer f.Close()

	pth, err := c.FileService.UploadFile(ctx, fh)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	var data interface{}
	if isDatapool {
		absPath := filepath.Join(dir.GetCurrentAbPath(), pth)
		data, err = c.DatapoolService.ReadExcel(absPath)

		if err != nil {
			ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
			return
		}
	}

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code,
		Data: iris.Map{"path": pth, "data": data}, Msg: _domain.NoErr.Msg})
}
