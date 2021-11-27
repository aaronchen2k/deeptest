package controller

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/aaronchen2k/deeptest/internal/server/core/web/validate"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
	"strings"

	"github.com/kataras/iris/v12"
	"go.uber.org/zap"
)

type ProductCtrl struct {
	ProductService *service.ProductService `inject:""`
	BaseCtrl
}

func NewProductCtrl() *ProductCtrl {
	return &ProductCtrl{}
}

// Query
// @summary 产品列表
// @Accept json
// @Produce json
// @Param req body domain.ProductReqPaginate true "Query Request Object"
// @Success 200 {object} domain.Response{data=domain.PageData} "code = success? 1 : 0"
// @Router /api/v1/products [get]
func (c *ProductCtrl) Query(ctx iris.Context) {
	var req serverDomain.ProductReqPaginate
	if err := ctx.ReadQuery(&req); err != nil {
		errs := validate.ValidRequest(err)
		if len(errs) > 0 {
			logUtils.Errorf("参数验证失败", zap.String("错误", strings.Join(errs, ";")))
			ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: strings.Join(errs, ";")})
			return
		}
	}
	req.ConvertParams()

	data, err := c.ProductService.Paginate(req)
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}
	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: data, Msg: domain.NoErr.Msg})
}

// Get 详情
func (c *ProductCtrl) Get(ctx iris.Context) {
	var req domain.ReqId
	if err := ctx.ReadParams(&req); err != nil {
		logUtils.Errorf("参数解析失败", zap.String("错误:", err.Error()))
		ctx.JSON(domain.Response{Code: domain.ParamErr.Code, Data: nil, Msg: domain.ParamErr.Msg})
		return
	}
	product, err := c.ProductService.FindById(req.Id)
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: domain.SystemErr.Msg})
		return
	}
	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: product, Msg: domain.NoErr.Msg})
}

// Create 添加
func (c *ProductCtrl) Create(ctx iris.Context) {
	req := serverDomain.ProductRequest{}
	if err := ctx.ReadJSON(&req); err != nil {
		errs := validate.ValidRequest(err)
		if len(errs) > 0 {
			logUtils.Errorf("参数验证失败", zap.String("错误", strings.Join(errs, ";")))
			ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: strings.Join(errs, ";")})
			return
		}
	}
	id, err := c.ProductService.Create(req)
	if err != nil {
		ctx.JSON(domain.Response{
			Code: c.ErrCode(err),
			Data: nil,
		})
		return
	}

	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: iris.Map{"id": id}, Msg: domain.NoErr.Msg})
}

// Update 更新
func (c *ProductCtrl) Update(ctx iris.Context) {
	var reqId domain.ReqId
	if err := ctx.ReadParams(&reqId); err != nil {
		logUtils.Errorf("参数解析失败", zap.String("错误:", err.Error()))
		ctx.JSON(domain.Response{Code: domain.ParamErr.Code, Data: nil, Msg: domain.ParamErr.Msg})
		return
	}

	var req serverDomain.ProductRequest
	if err := ctx.ReadJSON(&req); err != nil {
		errs := validate.ValidRequest(err)
		if len(errs) > 0 {
			logUtils.Errorf("参数验证失败", zap.String("错误", strings.Join(errs, ";")))
			ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: strings.Join(errs, ";")})
			return
		}
	}

	err := c.ProductService.Update(reqId.Id, req)
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}
	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: nil, Msg: domain.NoErr.Msg})
}

// Delete 删除
func (c *ProductCtrl) Delete(ctx iris.Context) {
	var req domain.ReqId
	if err := ctx.ReadParams(&req); err != nil {
		logUtils.Errorf("参数解析失败", zap.String("错误:", err.Error()))
		ctx.JSON(domain.Response{Code: domain.ParamErr.Code, Data: nil, Msg: domain.ParamErr.Msg})
		return
	}
	err := c.ProductService.DeleteById(req.Id)
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: nil, Msg: domain.NoErr.Msg})
}
