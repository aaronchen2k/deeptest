package handler

import (
	"github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/core/web/validate"
	"github.com/aaronchen2k/deeptest/internal/server/modules/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/repo"
	"github.com/aaronchen2k/deeptest/internal/server/modules/service"
	"github.com/aaronchen2k/deeptest/pkg/domain"
	"github.com/aaronchen2k/deeptest/pkg/lib/log"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/snowlyg/multi"
	"go.uber.org/zap"
)

type UserCtrl struct {
	UserService *service.UserService `inject:""`
	UserRepo    *repo.UserRepo       `inject:""`
}

func (c *UserCtrl) ListAll(ctx iris.Context) {
	var req domain.UserReqPaginate

	if err := ctx.ReadQuery(&req); err != nil {
		errs := validate.ValidRequest(err)
		if len(errs) > 0 {
			_logUtils.Errorf("参数验证失败", zap.String("错误", strings.Join(errs, ";")))
			ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: strings.Join(errs, ";")})
			return
		}
	}

	data, err := c.UserRepo.Paginate(req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: data, Msg: _domain.NoErr.Msg})
}

// GetUser 详情
func (c *UserCtrl) GetUser(ctx iris.Context) {
	var req _domain.ReqId
	if err := ctx.ReadParams(&req); err != nil {
		_logUtils.Errorf("参数解析失败", zap.String("错误:", err.Error()))
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: _domain.ParamErr.Msg})
		return
	}
	user, err := c.UserRepo.FindDetailById(req.Id)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: _domain.SystemErr.Msg})
		return
	}
	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: user, Msg: _domain.NoErr.Msg})
}

// Invite 邀请用户
func (c *UserCtrl) Invite(ctx iris.Context) {
	req := domain.InviteUserReq{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: err.Error()})
		return
	}

	_, bizErr := c.UserService.Invite(req)
	if bizErr != nil {
		ctx.JSON(_domain.Response{Code: bizErr.Code})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Msg: _domain.NoErr.Msg})
}

// UpdateEmail 修改邮箱
func (c *UserCtrl) UpdateEmail(ctx iris.Context) {
	userId := multi.GetUserId(ctx)
	req := domain.UpdateUserReq{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: err.Error()})
		return
	}

	po, _ := c.UserRepo.FindByEmail(req.Email, userId)
	if po.Id > 0 {
		bizErr := _domain.ErrEmailExist
		ctx.JSON(_domain.Response{Code: bizErr.Code})
		return
	}

	err = c.UserRepo.UpdateEmail(req.Email, userId)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	user, err := c.UserRepo.FindDetailById(userId)
	user.Password = ""
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: user, Msg: _domain.NoErr.Msg})
}

// UpdateName 修改名称
func (c *UserCtrl) UpdateName(ctx iris.Context) {
	userId := multi.GetUserId(ctx)
	req := domain.UpdateUserReq{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: err.Error()})
		return
	}

	po, _ := c.UserRepo.FindByUserName(req.Username, userId)
	if po.Id > 0 {
		bizErr := _domain.ErrUsernameExist
		ctx.JSON(_domain.Response{Code: bizErr.Code})
		return
	}

	err = c.UserRepo.UpdateName(req.Username, userId)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	user, err := c.UserRepo.FindDetailById(userId)
	user.Password = ""
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: user, Msg: _domain.NoErr.Msg})
}

// UpdatePassword 修改密码
func (c *UserCtrl) UpdatePassword(ctx iris.Context) {
	userId := multi.GetUserId(ctx)

	req := domain.UpdateUserReq{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: err.Error()})
		return
	}

	err = c.UserRepo.ChangePassword(req, userId)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	user, err := c.UserRepo.FindDetailById(userId)
	user.Password = ""
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: user, Msg: _domain.NoErr.Msg})
}

// Profile 个人信息
func (c *UserCtrl) Profile(ctx iris.Context) {
	id := multi.GetUserId(ctx)
	if id == 0 {
		ctx.JSON(_domain.Response{Code: _domain.ErrNoUser.Code, Msg: _domain.SystemErr.Msg})
		return
	}

	user, err := c.UserRepo.FindDetailById(id)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: _domain.SystemErr.Msg})
		return
	}

	user.Password = ""

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: user, Msg: _domain.NoErr.Msg})
}

// Message 消息
func (c *UserCtrl) Message(ctx iris.Context) {
	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Msg: _domain.NoErr.Msg})
}

// CreateUser 添加
func (c *UserCtrl) CreateUser(ctx iris.Context) {
	req := domain.UserReq{}
	if err := ctx.ReadJSON(&req); err != nil {
		errs := validate.ValidRequest(err)
		if len(errs) > 0 {
			_logUtils.Errorf("参数验证失败", zap.String("错误", strings.Join(errs, ";")))
			ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: strings.Join(errs, ";")})
			return
		}
	}
	id, err := c.UserRepo.Create(req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: iris.Map{"id": id}, Msg: _domain.NoErr.Msg})
}

// UpdateUser 更新
func (c *UserCtrl) UpdateUser(ctx iris.Context) {
	var reqId _domain.ReqId
	if err := ctx.ReadParams(&reqId); err != nil {
		_logUtils.Errorf("参数解析失败", zap.String("错误:", err.Error()))
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: _domain.ParamErr.Msg})
		return
	}

	var req domain.UserReq
	if err := ctx.ReadJSON(&req); err != nil {
		errs := validate.ValidRequest(err)
		if len(errs) > 0 {
			_logUtils.Errorf("参数验证失败", zap.String("错误", strings.Join(errs, ";")))
			ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: strings.Join(errs, ";")})
			return
		}
	}

	err := c.UserRepo.Update(reqId.Id, req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}
	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Msg: _domain.NoErr.Msg})
}

// DeleteUser 删除
func (c *UserCtrl) DeleteUser(ctx iris.Context) {
	var req _domain.ReqId
	if err := ctx.ReadParams(&req); err != nil {
		_logUtils.Errorf("参数解析失败", zap.String("错误:", err.Error()))
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: _domain.ParamErr.Msg})
		return
	}
	err := c.UserRepo.DeleteById(req.Id)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Msg: _domain.NoErr.Msg})
}

// Logout 退出
func (c *UserCtrl) Logout(ctx iris.Context) {
	token := multi.GetVerifiedToken(ctx)
	if token == nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: "授权凭证为空"})
		return
	}
	err := c.UserRepo.DelToken(string(token))
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}
	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Msg: _domain.NoErr.Msg})
}

// Clear 清空 token
func (c *UserCtrl) Clear(ctx iris.Context) {
	token := multi.GetVerifiedToken(ctx)
	if token == nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: "授权凭证为空"})
		return
	}
	if err := c.UserRepo.CleanToken(multi.AdminAuthority, string(token)); err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}
	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Msg: _domain.NoErr.Msg})
}

// ChangeAvatar 修改头像
func (c *UserCtrl) ChangeAvatar(ctx iris.Context) {
	avatar := &model.Avatar{}
	if err := ctx.ReadJSON(avatar); err != nil {
		errs := validate.ValidRequest(err)
		if len(errs) > 0 {
			_logUtils.Errorf("参数验证失败", zap.String("错误", strings.Join(errs, ";")))
			ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: strings.Join(errs, ";")})
			return
		}
	}
	err := c.UserRepo.UpdateAvatar(multi.GetUserId(ctx), avatar.Avatar)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}
	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Msg: _domain.NoErr.Msg})
}
