package handler

import (
	"github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/config"
	"github.com/aaronchen2k/deeptest/internal/server/core/web/validate"
	"github.com/aaronchen2k/deeptest/internal/server/modules/service"
	"github.com/aaronchen2k/deeptest/pkg/domain"
	logUtils "github.com/aaronchen2k/deeptest/pkg/lib/log"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/snowlyg/helper/str"
	"go.uber.org/zap"
)

type DataCtrl struct {
	DataService *service.DataService `inject:""`
}

// Init  初始化项目接口
// @summary 初始化项目
// @Tags 系统
// @Accept json
// @Produce	json
// @Router /api/v1/init/initdb [post]
// @Param ProjectReq body serverDomain.DataReq true "初始化db"
//@Success 200 {object} _domain.Response
func (c *DataCtrl) Init(ctx iris.Context) {
	req := serverDomain.DataReq{}
	if err := ctx.ReadJSON(&req); err != nil {
		errs := validate.ValidRequest(err)
		if len(errs) > 0 {
			logUtils.Errorf("参数验证失败", zap.String("错误", strings.Join(errs, ";")))
			ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: strings.Join(errs, ";")})
			return
		}
	}

	err := c.DataService.InitDB(req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: _domain.SystemErr.Msg})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Msg: _domain.NoErr.Msg})
}

// Check 检测是否需要初始化项目
func (c *DataCtrl) Check(ctx iris.Context) {
	if c.DataService.DataRepo.DB == nil {
		ctx.JSON(_domain.Response{Code: _domain.NeedInitErr.Code, Data: iris.Map{
			"needInit": true,
		}, Msg: str.Join(_domain.NeedInitErr.Msg, ":数据库初始化失败")})
		return
	} else if config.CONFIG.System.CacheType == "redis" && config.CACHE == nil {
		ctx.JSON(_domain.Response{Code: _domain.NeedInitErr.Code, Data: iris.Map{
			"needInit": true,
		}, Msg: str.Join(_domain.NeedInitErr.Msg, ":缓存驱动初始化失败")})
		return
	}
	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: iris.Map{
		"needInit": false,
	}, Msg: _domain.NoErr.Msg})
}
