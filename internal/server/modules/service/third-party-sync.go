package service

import (
	"encoding/json"
	"errors"
	"fmt"
	v1 "github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/config"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/core/cron"
	"github.com/aaronchen2k/deeptest/internal/pkg/helper/openapi/thirdPart"
	serverConsts "github.com/aaronchen2k/deeptest/internal/server/consts"
	"github.com/aaronchen2k/deeptest/internal/server/core/cache"
	"github.com/aaronchen2k/deeptest/internal/server/modules/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/repo"
	_commUtils "github.com/aaronchen2k/deeptest/pkg/lib/comm"
	logUtils "github.com/aaronchen2k/deeptest/pkg/lib/log"
	"github.com/getkin/kin-openapi/openapi3"
	"gorm.io/gorm"
	"time"
)

type ThirdPartySyncService struct {
	ThirdPartySyncRepo       *repo.ThirdPartySyncRepo    `inject:""`
	CategoryRepo             *repo.CategoryRepo          `inject:""`
	EndpointRepo             *repo.EndpointRepo          `inject:""`
	EndpointInterfaceRepo    *repo.EndpointInterfaceRepo `inject:""`
	UserRepo                 *repo.UserRepo              `inject:""`
	RemoteService            *RemoteService              `inject:""`
	ServeService             *ServeService               `inject:""`
	EndpointInterfaceService *EndpointInterfaceService   `inject:""`
	Cron                     *cron.ServerCron            `inject:""`
}

func (s *ThirdPartySyncService) GetToken(baseUrl string) (token string, err error) {
	loginByOauthReq := v1.LoginByOauthReq{
		LoginName: config.CONFIG.ThirdParty.Username,
		Password:  _commUtils.Sha256(config.CONFIG.ThirdParty.Password),
	}

	loginByOauthResData := s.RemoteService.LoginByOauth(loginByOauthReq, baseUrl)
	if loginByOauthResData.Code == "" {
		return "", errors.New("login fail")
	}

	getTokenFromCodeReq := v1.GetTokenFromCodeReq{
		Code: loginByOauthResData.Code,
	}

	getTokenFromCodeResData := s.RemoteService.GetTokenFromCode(getTokenFromCodeReq, baseUrl)
	token = getTokenFromCodeResData.Token

	return
}

func (s *ThirdPartySyncService) GetClasses(serviceCode, token string, baseUrl string) (classes []string) {
	findClassByServiceCodeReq := v1.FindClassByServiceCodeReq{
		ServiceCode: serviceCode,
	}
	findClassByServiceCodeResData := s.RemoteService.FindClassByServiceCode(findClassByServiceCodeReq, token, baseUrl)

	for _, v := range findClassByServiceCodeResData {
		classes = append(classes, v.Code)
	}

	return
}

func (s *ThirdPartySyncService) GetFunctionsByClass(serviceCode, classCode, token string, baseUrl string) (functions []string) {
	getFunctionsByClassReq := v1.GetFunctionsByClassReq{
		ServiceCode: serviceCode,
		ClassCode:   classCode,
	}
	getFunctionsByClassResData := s.RemoteService.GetFunctionsByClass(getFunctionsByClassReq, token, baseUrl)
	for _, v := range getFunctionsByClassResData {
		//不同步内部方法
		if v.MessageType == 1 {
			functions = append(functions, v.Code)
		}
	}

	return
}

func (s *ThirdPartySyncService) GetFunctionDetail(classCode, function, token string, baseUrl string) (data v1.MetaGetMethodDetailResData) {
	metaGetMethodDetailReq := v1.MetaGetMethodDetailReq{
		ClassName:   classCode,
		Method:      function,
		IncludeSelf: true,
	}
	data = s.RemoteService.MetaGetMethodDetail(metaGetMethodDetailReq, token, baseUrl)

	return
}

func (s *ThirdPartySyncService) SaveData() (err error) {
	thirdPartySyncStatus, _ := cache.GetCacheString("thirdPartySyncStatus")
	if thirdPartySyncStatus == "Start" {
		return
	}

	_ = cache.SetCache("thirdPartySyncStatus", "Start", 4*time.Hour)
	syncList, err := s.GetAllData()
	if err != nil {
		return
	}

	for _, syncConfig := range syncList {
		if syncConfig.Switch == consts.SwitchOFF {
			continue
		}
		projectId, serveId, userId, baseUrl := syncConfig.ProjectId, syncConfig.ServeId, syncConfig.CreateUserId, syncConfig.Url

		token, err := s.GetToken(baseUrl)
		if err != nil {
			continue
		}

		classes := s.GetClasses(syncConfig.ServiceCode, token, baseUrl)
		for _, class := range classes {
			categoryId, err := s.SaveCategory(class, projectId, syncConfig.ServeId)
			if err != nil {
				continue
			}

			functionList := s.GetFunctionsByClass(syncConfig.ServiceCode, class, token, baseUrl)
			for _, function := range functionList {
				path := "/" + syncConfig.ServiceCode + "/" + class + "/" + function
				functionDetail := s.GetFunctionDetail(class, function, token, baseUrl)
				if functionDetail.Code == "" {
					continue
				}

				title := class + "-" + functionDetail.Code
				endpoint, err := s.EndpointRepo.GetByItem(consts.ThirdPartySync, projectId, path, serveId, title)
				if err != nil && err != gorm.ErrRecordNotFound {
					continue
				}

				var endpointId uint
				//if endpoint.ID != 0 && syncType != consts.Add {
				if endpoint.ID != 0 {
					endpointId = endpoint.ID
				}

				endpointId, err = s.SaveEndpoint(title, projectId, serveId, userId, endpointId, int64(categoryId), path)
				if err != nil {
					continue
				}

				interfaceId, err := s.SaveEndpointInterface(title, functionDetail, endpointId, projectId, path)
				if err != nil {
					continue
				}

				if err = s.SaveBody(functionDetail, interfaceId); err != nil {
					continue
				}
			}
		}
	}

	cache.SetCache("thirdPartySyncStatus", "Stop", -1)
	return
}

func (s *ThirdPartySyncService) SaveCategory(classCode string, projectId, serveId uint) (categoryId uint, err error) {
	rootNode, err := s.CategoryRepo.GetRootNode(projectId, serverConsts.EndpointCategory)
	if err != nil {
		return
	}

	categoryReq := model.Category{
		Name:       classCode,
		ProjectId:  projectId,
		ServeId:    serveId,
		Type:       serverConsts.EndpointCategory,
		SourceType: consts.ThirdPartySync,
		ParentId:   int(rootNode.ID),
	}

	category, err := s.CategoryRepo.GetDetail(categoryReq)
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}

	if category.ID != 0 {
		return category.ID, nil
	}

	err = s.CategoryRepo.Save(&categoryReq)
	if err != nil {
		return 0, err
	}
	categoryId = categoryReq.ID

	return
}

func (s *ThirdPartySyncService) SaveEndpoint(title string, projectId, serveId, userId, oldEndpointId uint, categoryId int64, path string) (endpointId uint, err error) {
	endpoint := model.Endpoint{
		Title:      title,
		ProjectId:  projectId,
		ServeId:    serveId,
		Path:       path,
		Status:     1,
		CategoryId: categoryId,
		SourceType: consts.ThirdPartySync,
	}

	if userId != 0 {
		user, err := s.UserRepo.FindById(userId)
		if err != nil {
			return 0, err
		}
		endpoint.CreateUser = user.Username
	}

	if oldEndpointId != 0 {
		endpoint.ID = oldEndpointId
	}

	err = s.EndpointRepo.SaveAll(&endpoint)
	if err != nil {
		return 0, err
	}

	endpointId = endpoint.ID

	return
}

func (s *ThirdPartySyncService) SaveEndpointInterface(title string, functionDetail v1.MetaGetMethodDetailResData, endpointId, projectId uint, path string) (interfaceId uint, err error) {
	endpointInterface := model.EndpointInterface{
		InterfaceBase: model.InterfaceBase{
			Name:      title,
			ProjectId: projectId,
			InterfaceConfigBase: model.InterfaceConfigBase{
				Url:      path,
				Method:   "POST",
				BodyType: s.getBodyType(functionDetail.RequestType),
				Version:  "v0.1.0",
			},
		},
		EndpointId:    endpointId,
		SourceType:    consts.ThirdPartySync,
		ResponseCodes: "200",
	}
	err = s.EndpointInterfaceRepo.Save(&endpointInterface)

	interfaceId = endpointInterface.ID
	return
}

func (s *ThirdPartySyncService) getBodyType(requestType string) (bodyType consts.HttpContentType) {
	if requestType == "JSON" {
		bodyType = consts.ContentTypeJSON
	} else if requestType == "FORM" {
		bodyType = consts.ContentTypeFormData
	}

	return
}

func (s *ThirdPartySyncService) getRequestBody(functionDetail v1.MetaGetMethodDetailResData) (requestBody string) {
	if functionDetail.RequestType == "JSON" {
		requestBody = functionDetail.RequestBody
	} else if functionDetail.RequestType == "FORM" {
		requestBody = functionDetail.RequestFormBody
	}

	return
}

func (s *ThirdPartySyncService) GetSchema(bodyString, requestType string) (schema *openapi3.SchemaRef) {
	if bodyString == "" {
		return
	}

	var schemas thirdPart.Schemas
	_ = json.Unmarshal([]byte(bodyString), &schemas)

	if requestType == "JSON" {
		schemas = schemas["root"].Properties
	}

	return thirdPart.NewThirdPart2conv().Convert(schemas)

}

func (s *ThirdPartySyncService) SaveBody(functionDetail v1.MetaGetMethodDetailResData, interfaceId uint) (err error) {
	functionBody := functionDetail.RequestBody
	if functionDetail.RequestType == "FORM" {
		functionBody = functionDetail.RequestFormBody
	}

	requestBodySchema := s.GetSchema(functionBody, functionDetail.RequestType)
	responseBodySchema := s.GetSchema(functionDetail.ResponseBody, functionDetail.RequestType)
	fmt.Println(requestBodySchema, responseBodySchema)

	//requestSchema := s.ServeService.Example2Schema(s.getRequestBody(functionDetail))
	requestSchemaString, _ := json.Marshal(requestBodySchema)

	generateFromRequestReq := v1.GenerateFromRequestReq{
		ContentType: s.getBodyType(functionDetail.RequestType).String(),
		Data:        string(requestSchemaString),
		InterfaceId: interfaceId,
	}
	_, err = s.EndpointInterfaceService.GenerateFromRequest(generateFromRequestReq)
	if err != nil {
		return
	}

	//responseSchema := s.ServeService.Example2Schema(functionDetail.ResponseBody)
	responseSchemaString, _ := json.Marshal(responseBodySchema)

	generateFromResponseReq := v1.GenerateFromResponseReq{
		Code:        "200",
		ContentType: s.getBodyType(functionDetail.ResponseType).String(),
		Data:        string(responseSchemaString),
		InterfaceId: interfaceId,
	}
	_, err = s.EndpointInterfaceService.GenerateFromResponse(generateFromResponseReq)
	return
}

func (s *ThirdPartySyncService) UpdateExecTimeById(id uint) (err error) {
	return s.ThirdPartySyncRepo.UpdateExecTimeById(id)
}

func (s *ThirdPartySyncService) AddThirdPartySyncCron() {
	name := "ThirdPartySync"

	s.Cron.RemoveTask(name)

	s.Cron.AddCommonTask(name, "* */12 * * *", func() {
		err := s.SaveData()
		if err != nil {
			logUtils.Error("third party 定时导入任务失败，错误原因：" + err.Error())
		}

		logUtils.Info("third party 定时任务结束")
	})
}

func (s *ThirdPartySyncService) GetAllData() (res []model.ThirdPartySync, err error) {
	return s.ThirdPartySyncRepo.AllData()
}