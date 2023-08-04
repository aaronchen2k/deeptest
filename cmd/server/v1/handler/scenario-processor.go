package handler

import (
	domain "github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/server/modules/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/service"
	"github.com/aaronchen2k/deeptest/pkg/domain"
	logUtils "github.com/aaronchen2k/deeptest/pkg/lib/log"
	"github.com/kataras/iris/v12"
)

type ScenarioProcessorCtrl struct {
	ScenarioProcessorService *service.ScenarioProcessorService `inject:""`
	BaseCtrl
}

// Get 详情
// @Tags	场景模块/处理器
// @summary	场景节点详情
// @accept 	application/json
// @Produce application/json
// @Param	Authorization	header	string	true	"Authentication header"
// @Param 	currProjectId	query	int		true	"当前项目ID"
// @Param 	id				path	int		true	"节点ID"
// @success	200	{object}	_domain.Response
// @Router	/api/v1/scenarios/processors/{id}	[get]
func (c *ScenarioProcessorCtrl) Get(ctx iris.Context) {
	processorId, err := ctx.Params().GetInt("id")
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: _domain.ParamErr.Msg})
		return
	}

	processorEntity, err := c.ScenarioProcessorService.GetEntity(processorId)

	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: _domain.SystemErr.Msg})
		return
	}
	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: processorEntity})
}

// UpdateName 更新
//func (c *ScenarioProcessorCtrl) UpdateName(ctx iris.Context) {
//	var req agentExec.ProcessorEntityBase
//	err := ctx.ReadJSON(&req)
//	if err != nil {
//		logUtils.Errorf("参数验证失败", err.Error())
//		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: _domain.ParamErr.Msg})
//		return
//	}
//
//	err = c.ScenarioProcessorService.UpdateName(req)
//	if err != nil {
//		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
//		return
//	}
//
//	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Msg: _domain.NoErr.Msg})
//}

// SaveBasicInfo 更新
// @Tags	场景模块/处理器
// @summary	保存基本信息
// @accept 	application/json
// @Produce application/json
// @Param	Authorization			header	string							true	"Authentication header"
// @Param 	currProjectId			query	int								true	"当前项目ID"
// @Param 	ScenarioProcessorInfo	body	domain.ScenarioProcessorInfo	true	"保存基本信息的请求参数"
// @success	200	{object}	_domain.Response{data=object{name=string}}
// @Router	/api/v1/scenarios/processors/saveProcessorInfo	[put]
func (c *ScenarioProcessorCtrl) SaveBasicInfo(ctx iris.Context) {
	var req domain.ScenarioProcessorInfo
	err := ctx.ReadJSON(&req)
	if err != nil {
		logUtils.Errorf("参数验证失败", err.Error())
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: _domain.ParamErr.Msg})
		return
	}

	err = c.ScenarioProcessorService.SaveProcessorInfo(req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: req.Name})
}

// Save 保存
// @Tags	场景模块/处理器
// @summary	保存配置信息
// @accept 	application/json
// @Produce application/json
// @Param	Authorization			header	string							true	"Authentication header"
// @Param 	currProjectId			query	int								true	"当前项目ID"
// @Param 	category			path	string								true	"category"
// @success	200	{object}	_domain.Response{data=model.ProcessorData}
// @Router	/api/v1/scenarios/processors/{category}/save	[put]
func (c *ScenarioProcessorCtrl) Save(ctx iris.Context) {
	processorCategoryString := ctx.Params().Get("category")
	processorCategory := consts.ProcessorCategory(processorCategoryString)

	var err error
	var po interface{}

	if processorCategory == consts.ProcessorGroup {
		var entity model.ProcessorGroup
		err = ctx.ReadJSON(&entity)
		err = c.ScenarioProcessorService.SaveGroup(&entity)
		po = entity

	} else if processorCategory == consts.ProcessorLogic {
		var entity model.ProcessorLogic
		err = ctx.ReadJSON(&entity)
		err = c.ScenarioProcessorService.SaveLogic(&entity)
		po = entity

	} else if processorCategory == consts.ProcessorLoop {
		var entity model.ProcessorLoop
		err = ctx.ReadJSON(&entity)
		err = c.ScenarioProcessorService.SaveLoop(&entity)
		po = entity

	} else if processorCategory == consts.ProcessorTimer {
		var entity model.ProcessorTimer
		err = ctx.ReadJSON(&entity)
		err = c.ScenarioProcessorService.SaveTimer(&entity)
		po = entity

	} else if processorCategory == consts.ProcessorPrint {
		var entity model.ProcessorPrint
		err = ctx.ReadJSON(&entity)
		err = c.ScenarioProcessorService.SavePrint(&entity)
		po = entity

	} else if processorCategory == consts.ProcessorVariable {
		var entity model.ProcessorVariable
		err = ctx.ReadJSON(&entity)
		err = c.ScenarioProcessorService.SaveVariable(&entity)
		po = entity

	} else if processorCategory == consts.ProcessorCookie {
		var entity model.ProcessorCookie
		err = ctx.ReadJSON(&entity)
		err = c.ScenarioProcessorService.SaveCookie(&entity)
		po = entity

	} else if processorCategory == consts.ProcessorAssertion {
		var entity model.ProcessorAssertion
		err = ctx.ReadJSON(&entity)
		err = c.ScenarioProcessorService.SaveAssertion(&entity)
		po = entity

	} else if processorCategory == consts.ProcessorExtractor {
		var entity model.ProcessorExtractor
		err = ctx.ReadJSON(&entity)
		err = c.ScenarioProcessorService.SaveExtractor(&entity)
		po = entity

	} else if processorCategory == consts.ProcessorData {
		var entity model.ProcessorData
		err = ctx.ReadJSON(&entity)
		err = c.ScenarioProcessorService.SaveData(&entity)
		po = entity

	}

	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: po, Msg: _domain.NoErr.Msg})
}
