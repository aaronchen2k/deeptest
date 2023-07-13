package service

import (
	"fmt"
	v1 "github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/helper/openapi"
	"github.com/aaronchen2k/deeptest/internal/server/modules/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/repo"
	_domain "github.com/aaronchen2k/deeptest/pkg/domain"
	_commUtils "github.com/aaronchen2k/deeptest/pkg/lib/comm"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/jinzhu/copier"
	encoder "github.com/zwgblue/yaml-encoder"
)

type ServeService struct {
	ServeRepo       *repo.ServeRepo       `inject:""`
	ServeServerRepo *repo.ServeServerRepo `inject:""`

	EndpointRepo          *repo.EndpointRepo          `inject:""`
	EndpointInterfaceRepo *repo.EndpointInterfaceRepo `inject:""`
}

func (s *ServeService) ListByProject(projectId int, userId uint) (ret []model.Serve, currServe model.Serve, err error) {
	ret, err = s.ServeRepo.ListByProject(projectId)

	currServe, err = s.ServeRepo.GetCurrServeByUser(userId)

	if currServe.ProjectId != uint(projectId) { //重新更新默认服务
		if len(ret) > 0 {
			currServe, err = s.ServeRepo.ChangeServe(ret[0].ID, userId)
		}
	}

	return
}

func (s *ServeService) Paginate(req v1.ServeReqPaginate) (ret _domain.PageData, err error) {
	ret, err = s.ServeRepo.Paginate(req)
	return
}

func (s *ServeService) Save(req v1.ServeReq) (res uint, err error) {
	var serve model.Serve
	if s.ServeRepo.ServeExist(uint(req.ID), req.ProjectId, req.Name) {
		err = fmt.Errorf("serve name already exist")
		return
	}
	copier.CopyWithOption(&serve, req, copier.Option{DeepCopy: true})
	err = s.ServeRepo.SaveServe(&serve)
	return serve.ID, err
}

func (s *ServeService) GetById(id uint) (res model.Serve) {
	res, _ = s.ServeRepo.Get(id)
	return
}

func (s *ServeService) DeleteById(id uint) (err error) {
	/*
		err = s.canDelete(id)
		if err != nil {
			return err
		}
	*/
	err = s.ServeRepo.DeleteById(id)
	return
}

func (s *ServeService) canDelete(id uint) (err error) {
	var count int64
	count, err = s.EndpointRepo.GetCountByServeId(id)
	if err != nil {
		return
	}
	if count != 0 {
		err = fmt.Errorf("interfaces are created under the service,not allowed to delete")
	}

	return
}

func (s *ServeService) DisableById(id uint) (err error) {
	err = s.ServeRepo.DisableById(id)
	return
}

func (s *ServeService) PaginateVersion(req v1.ServeVersionPaginate) (ret _domain.PageData, err error) {
	return s.ServeRepo.PaginateVersion(req)
}

func (s *ServeService) SaveVersion(req v1.ServeVersionReq) (res uint, err error) {
	var serveVersion model.ServeVersion
	if s.ServeRepo.VersionExist(uint(req.ID), uint(req.ServeId), req.Value) {
		err = fmt.Errorf("serve version already exist")
		return
	}
	copier.CopyWithOption(&serveVersion, req, copier.Option{DeepCopy: true})
	err, res = s.ServeRepo.SaveVersion(serveVersion.ID, &serveVersion), serveVersion.ID
	return
}

func (s *ServeService) DeleteVersionById(id uint) (err error) {
	err = s.ServeRepo.DeleteVersionById(id)
	return
}

func (s *ServeService) DisableVersionById(id uint) (err error) {
	err = s.ServeRepo.DisableVersionById(id)
	return
}

func (s *ServeService) ListServer(req v1.ServeServer) (res []model.ServeServer, err error) {
	if req.ServeId == 0 {
		server, _ := s.ServeServerRepo.Get(req.ServerId)
		req.ServeId = server.ServeId
	}

	res, err = s.ServeRepo.ListServer(req.ServeId)

	return
}

func (s *ServeService) SaveServer(req v1.ServeServer) (res uint, err error) {
	var serve model.ServeServer
	copier.CopyWithOption(&serve, req, copier.Option{DeepCopy: true})
	err = s.ServeRepo.Save(serve.ID, &serve)
	return serve.ID, err
}

func (s *ServeService) Copy(id uint) (err error) {
	serve, _ := s.ServeRepo.Get(id)
	serve.ID = 0
	serve.Name += "_copy"
	serve.CreatedAt = nil
	serve.UpdatedAt = nil
	return s.ServeRepo.Save(0, &serve)
}

func (s *ServeService) SaveSchema(req v1.ServeSchemaReq) (res uint, err error) {
	var serveSchema model.ComponentSchema
	if req.ID == 0 && s.ServeRepo.SchemaExist(uint(req.ID), uint(req.ServeId), req.Name) {
		err = fmt.Errorf("schema name already exist")
		return
	}
	copier.CopyWithOption(&serveSchema, req, copier.Option{DeepCopy: true})
	serveSchema.Ref = "#/components/schemas/" + serveSchema.Name
	err = s.ServeRepo.Save(serveSchema.ID, &serveSchema)
	return serveSchema.ID, err
}

func (s *ServeService) SaveSecurity(req v1.ServeSecurityReq) (res uint, err error) {
	var serveSecurity model.ComponentSchemaSecurity
	if s.ServeRepo.SecurityExist(uint(req.ID), uint(req.ServeId), req.Name) {
		err = fmt.Errorf("security name already exist")
		return
	}
	copier.CopyWithOption(&serveSecurity, req, copier.Option{DeepCopy: true})
	err = s.ServeRepo.Save(serveSecurity.ID, &serveSecurity)
	return serveSecurity.ID, err
}

func (s *ServeService) PaginateSchema(req v1.ServeSchemaPaginate) (ret _domain.PageData, err error) {
	return s.ServeRepo.PaginateSchema(req)
}

func (s *ServeService) GetSchema(serverId uint, ref string) (schema model.ComponentSchema, err error) {
	return s.ServeRepo.GetSchemaByRef(serverId, ref)
}

func (s *ServeService) PaginateSecurity(req v1.ServeSecurityPaginate) (ret _domain.PageData, err error) {
	return s.ServeRepo.PaginateSecurity(req)
}

func (s *ServeService) Example2Schema(data string) (schema openapi.Schema) {
	schema2conv := openapi.NewSchema2conv()
	var obj interface{}
	schema = openapi.Schema{}
	_commUtils.JsonDecode(data, &obj)
	//_commUtils.JsonDecode("{\"id\":1,\"name\":\"user\"}", &obj)
	//_commUtils.JsonDecode("[\"0，2，3\"]", &obj)
	//_commUtils.JsonDecode("[]", &obj)
	//_commUtils.JsonDecode("[{\"id\":1,\"name\":\"user\"}]", &obj)
	//_commUtils.JsonDecode("{\"id\":[1,2,3],\"name\":\"user\"}", &obj)
	schema2conv.Example2Schema(obj, &schema)
	//fmt.Println(_commUtils.JsonEncode(schema), "++++++++++++")
	return
}

func (s *ServeService) DeleteSchemaById(id uint) (err error) {
	//TODO
	var schema model.ComponentSchema
	schema, err = s.ServeRepo.GetSchema(id)
	if err != nil {
		return err
	}

	var count int64
	count, err = s.EndpointInterfaceRepo.GetCountByRef(schema.Ref)
	if err != nil {
		return
	}

	if count > 0 {
		err = fmt.Errorf("the schema has been referenced and cannot be deleted")
		return
	}

	err = s.ServeRepo.DeleteSchemaById(id)
	return
}

func (s *ServeService) DeleteSecurityId(id uint) (err error) {
	err = s.ServeRepo.DeleteSecurityId(id)
	return
}

func (s *ServeService) Schema2Example(serveId uint, data string) (obj interface{}) {
	schema2conv := openapi.NewSchema2conv()
	schema2conv.Components = s.Components(serveId)
	//schema1 := openapi3.Schema{}
	//_commUtils.JsonDecode(data, &schema)
	//_commUtils.JsonDecode("{\"type\":\"array\",\"items\":{\"type\":\"number\"}}", &schema)
	//_commUtils.JsonDecode("{\"properties\":{\"id\":{\"type\":\"number\"},\"name\":{\"type\":\"string\"}},\"type\":\"object\"}", &schema)
	//_commUtils.JsonDecode("{\"type\":\"array\",\"items\":{\"properties\":{\"id\":{\"type\":\"number\"},\"name\":{\"type\":\"string\"}},\"type\":\"object\"}}", &schema)
	schema := openapi.SchemaRef{}
	//data = "{\"type\":\"object\",\"properties\":{\"name1\":{\"type\":\"object\",\"ref\":\"#/components/schemas/user1\",\"name\":\"user1\"},\"name2\":{\"type\":\"string\"},\"name3\":{\"type\":\"string\"}}}"
	_commUtils.JsonDecode(data, &schema)
	//_commUtils.JsonDecode("{\"type\":\"array\",\"items\":{\"type\":\"number\"}}", &schema1)
	//copier.CopyWithOption(&schema, a, copier.Option{DeepCopy: true})
	//fmt.Println(schema, "+++++++++++++")
	obj = schema2conv.Schema2Example(schema)
	//fmt.Println(schema.Items, "+++++", schema1.Items, _commUtils.JsonEncode(obj), "++++++++++++")
	return
}

func (s *ServeService) Components(serveId uint) (components openapi.Components) {
	components = openapi.Components{}
	result, err := s.ServeRepo.GetSchemasByServeId(serveId)
	if err != nil {
		return
	}

	for _, item := range result {
		var schema openapi.SchemaRef
		_commUtils.JsonDecode(item.Content, &schema)
		components[item.Ref] = schema
	}

	return

}

func (s *ServeService) Schema2Yaml(data string) (res string) {
	schema := openapi3.Schema{}
	_commUtils.JsonDecode(data, &schema)
	content, _ := encoder.NewEncoder(schema).Encode()
	return string(content)
}

func (s *ServeService) CopySchema(id uint) (schema model.ComponentSchema, err error) {
	schema, err = s.ServeRepo.GetSchema(id)
	if err != nil {
		return
	}
	schema.ID = 0
	schema.CreatedAt = nil
	schema.UpdatedAt = nil
	err = s.ServeRepo.Save(0, &schema)
	return
}

func (s *ServeService) BindEndpoint(req v1.ServeVersionBindEndpointReq) (err error) {
	var serveEndpointVersion []model.ServeEndpointVersion
	for _, endpointVersion := range req.EndpointVersions {
		serveEndpointVersion = append(serveEndpointVersion, model.ServeEndpointVersion{EndpointId: endpointVersion.EndpointId, EndpointVersion: endpointVersion.Version, ServeId: req.ServeId, ServeVersion: req.ServeVersion})
	}
	err = s.ServeRepo.BindEndpoint(req.ServeId, req.ServeVersion, serveEndpointVersion)
	return
}

func (s *ServeService) ChangeServe(serveId, userId uint) (serve model.Serve, err error) {
	serve, err = s.ServeRepo.ChangeServe(serveId, userId)

	return
}

func (s *ServeService) SaveSwaggerSync(req v1.SwaggerSyncReq) (data model.SwaggerSync, err error) {
	var swaggerSync model.SwaggerSync
	copier.CopyWithOption(&swaggerSync, req, copier.Option{DeepCopy: true})
	err = s.ServeRepo.SaveSwaggerSync(&swaggerSync)
	data = swaggerSync
	return
}

func (s *ServeService) SwaggerSyncDetail(projectId uint) (data model.SwaggerSync, err error) {
	return s.ServeRepo.GetSwaggerSync(projectId)
}

func (s *ServeService) SwaggerSyncList() (data []model.SwaggerSync, err error) {
	return s.ServeRepo.GetSwaggerSyncList()
}
