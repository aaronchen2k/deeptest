package handler

import (
	scriptHelper "github.com/aaronchen2k/deeptest/internal/pkg/helper/script"
	service "github.com/aaronchen2k/deeptest/internal/server/modules/service"
	"github.com/aaronchen2k/deeptest/pkg/domain"
	"github.com/kataras/iris/v12"
)

type SnippetCtrl struct {
	SnippetService *service.SnippetService `inject:""`

	BaseCtrl
}

// Get 详情
// @Tags	脚本
// @summary	获取详情
// @accept 	application/json
// @Produce application/json
// @Param	Authorization	header	string	true	"Authentication header"
// @Param 	currProjectId	query	int		true	"当前项目ID"
// @Param 	name			query	string	true	"脚本名"
// @success	200	{object}	_domain.Response{data=model.Snippet}
// @Router	/api/v1/snippets	[get]
func (c *SnippetCtrl) Get(ctx iris.Context) {
	name := ctx.URLParam("name")

	snippet, err := c.SnippetService.Get(scriptHelper.ScriptType(name))
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: _domain.SystemErr.Msg})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: snippet})
}
