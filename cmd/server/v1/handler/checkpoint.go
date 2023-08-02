package handler

import (
	"github.com/aaronchen2k/deeptest/internal/server/modules/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/service"
	"github.com/aaronchen2k/deeptest/pkg/domain"
	"github.com/kataras/iris/v12"
)

type CheckpointCtrl struct {
	CheckpointService *service.CheckpointService `inject:""`
	BaseCtrl
}

// Get 详情
// @Tags	检查点
// @summary	检查点详情
// @accept 	application/json
// @Produce application/json
// @Param	Authorization		header	string	true	"Authentication header"
// @Param 	currProjectId		query	int		true	"当前项目ID"
// @Param 	id					path	int		true	"检查点ID"
// @success	200	{object}	_domain.Response{data=model.DebugInterfaceCheckpoint}
// @Router	/api/v1/checkpoints/{id}	[get]
func (c *CheckpointCtrl) Get(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: _domain.ParamErr.Msg})
		return
	}

	checkpoint, err := c.CheckpointService.Get(uint(id))
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: _domain.SystemErr.Msg})
		return
	}
	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: checkpoint})
}

// Create 添加
// @Tags	检查点
// @summary	新建检查点
// @accept 	application/json
// @Produce application/json
// @Param	Authorization				header	string								true	"Authentication header"
// @Param 	currProjectId				query	int									true	"当前项目ID"
// @Param 	DebugInterfaceCheckpoint	body	model.DebugInterfaceCheckpoint		true	"新建检查点的请求体"
// @success	200	{object}	_domain.Response{data=model.DebugInterfaceCheckpoint}
// @Router	/api/v1/checkpoints	[post]
func (c *CheckpointCtrl) Create(ctx iris.Context) {
	checkpoint := model.DebugConditionCheckpoint{}
	err := ctx.ReadJSON(&checkpoint)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: _domain.ParamErr.Msg})
		return
	}

	err = c.CheckpointService.Create(&checkpoint)
	if err != nil {
		ctx.JSON(_domain.Response{
			Code: _domain.SystemErr.Code,
		})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: checkpoint, Msg: _domain.NoErr.Msg})
}

// Update 更新
// @Tags	检查点
// @summary	更新检查点
// @accept 	application/json
// @Produce application/json
// @Param	Authorization				header	string								true	"Authentication header"
// @Param 	currProjectId				query	int									true	"当前项目ID"
// @Param 	DebugInterfaceCheckpoint	body	model.DebugInterfaceCheckpoint		true	"更新检查点的请求体"
// @success	200	{object}	_domain.Response
// @Router	/api/v1/checkpoints	[put]
func (c *CheckpointCtrl) Update(ctx iris.Context) {
	var checkpoint model.DebugConditionCheckpoint
	err := ctx.ReadJSON(&checkpoint)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: _domain.ParamErr.Msg})
		return
	}

	err = c.CheckpointService.Update(&checkpoint)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}
	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Msg: _domain.NoErr.Msg})
}

// Delete 删除
// @Tags	检查点
// @summary	删除检查点
// @accept 	application/json
// @Produce application/json
// @Param	Authorization		header	string	true	"Authentication header"
// @Param 	currProjectId		query	int		true	"当前项目ID"
// @Param 	id					path	int		true	"检查点ID"
// @success	200	{object}	_domain.Response
// @Router	/api/v1/checkpoints/{id}	[delete]
func (c *CheckpointCtrl) Delete(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: _domain.ParamErr.Msg})
		return
	}

	err = c.CheckpointService.Delete(uint(id))
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Msg: _domain.NoErr.Msg})
}
