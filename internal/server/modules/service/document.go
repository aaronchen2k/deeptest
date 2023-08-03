package service

import (
	"encoding/base64"
	domain "github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/repo"
	commUtils "github.com/aaronchen2k/deeptest/pkg/lib/comm"
	"github.com/jinzhu/copier"
	"strconv"
	"strings"
)

type DocumentService struct {
	EndpointRepo          *repo.EndpointRepo          `inject:""`
	ProjectRepo           *repo.ProjectRepo           `inject:""`
	ServeRepo             *repo.ServeRepo             `inject:""`
	EnvironmentRepo       *repo.EnvironmentRepo       `inject:""`
	EndpointDocumentRepo  *repo.EndpointDocumentRepo  `inject:""`
	EndpointSnapshotRepo  *repo.EndpointSnapshotRepo  `inject:""`
	EndpointInterfaceRepo *repo.EndpointInterfaceRepo `inject:""`
	EndpointService       *EndpointService            `inject:""`
}

const (
	EncryptKey = "docencryptkey123"
)

func (s *DocumentService) Content(req domain.DocumentReq) (res domain.DocumentRep, err error) {
	var projectId, documentId uint
	var endpointIds, serveIds []uint
	var needDetail bool

	projectId, serveIds, endpointIds, documentId, needDetail = req.ProjectId, req.ServeIds, req.EndpointIds, req.DocumentId, req.NeedDetail

	var endpoints map[uint][]domain.EndpointReq
	endpoints, err = s.GetEndpoints(&projectId, &serveIds, &endpointIds, documentId, needDetail)
	if err != nil {
		return
	}

	res = s.GetProject(projectId)

	res.Serves = s.GetServes(serveIds, endpoints)

	return
}

func (s *DocumentService) GetEndpoints(projectId *uint, serveIds, endpointIds *[]uint, documentId uint, needDetail bool) (res map[uint][]domain.EndpointReq, err error) {
	var endpoints []*model.Endpoint

	if documentId != 0 {
		endpoints, err = s.EndpointSnapshotRepo.GetByDocumentId(documentId)
	} else if *projectId != 0 {
		endpoints, err = s.EndpointRepo.GetByProjectId(*projectId, needDetail)
	} else if len(*serveIds) != 0 {
		endpoints, err = s.EndpointRepo.GetByServeIds(*serveIds, needDetail)
	} else if len(*endpointIds) != 0 {
		endpoints, err = s.EndpointRepo.GetByEndpointIds(*endpointIds, needDetail)
	}

	if err != nil {
		return
	}

	res = s.GetEndpointsInfo(projectId, serveIds, endpoints)

	return
}

func (s *DocumentService) GetEndpointsInfo(projectId *uint, serveIds *[]uint, endpoints []*model.Endpoint) (res map[uint][]domain.EndpointReq) {
	res = make(map[uint][]domain.EndpointReq)

	serves := make(map[uint]uint)
	for _, item := range endpoints {
		var endpoint domain.EndpointReq
		//ret, _ := s.EndpointRepo.GetAll(item.ID, "v0.1.0")
		copier.CopyWithOption(&endpoint, &item, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		res[endpoint.ServeId] = append(res[endpoint.ServeId], endpoint)
		if _, ok := serves[endpoint.ServeId]; !ok {
			*serveIds = append(*serveIds, endpoint.ServeId)
			serves[endpoint.ServeId] = endpoint.ServeId
		}

		*projectId = endpoint.ProjectId
	}
	return
}

func (s *DocumentService) GetProject(projectId uint) (doc domain.DocumentRep) {
	project, err := s.ProjectRepo.Get(projectId)
	if err != nil {
		return
	}
	copier.CopyWithOption(&doc, &project, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	doc.GlobalParams, _ = s.EnvironmentRepo.ListParams(projectId)
	doc.GlobalVars = s.GetGlobalVars(projectId)
	return
}

func (s *DocumentService) GetServes(serveIds []uint, endpoints map[uint][]domain.EndpointReq) (serves []domain.DocumentServe) {
	res, _ := s.ServeRepo.GetServesByIds(serveIds)
	schemas := s.GetSchemas(serveIds)
	securities := s.GetSecurities(serveIds)
	servers := s.GetServers(serveIds)
	for _, item := range res {
		var serve domain.DocumentServe
		copier.CopyWithOption(&serve, &item, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		serve.Endpoints = endpoints[uint(serve.ID)]
		serve.Component = schemas[uint(serve.ID)]
		serve.Securities = securities[uint(serve.ID)]
		serve.Servers = servers[uint(serve.ID)]
		serves = append(serves, serve)
	}
	return
}

func (s *DocumentService) GetSchemas(serveIds []uint) (schemas map[uint][]domain.ServeSchemaReq) {
	schemas = make(map[uint][]domain.ServeSchemaReq)
	res, _ := s.ServeRepo.GetSchemas(serveIds)
	for _, item := range res {
		var schema domain.ServeSchemaReq
		copier.CopyWithOption(&schema, &item, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		schemas[uint(schema.ServeId)] = append(schemas[uint(schema.ServeId)], schema)
	}
	return
}

func (s *DocumentService) GetServers(serveIds []uint) (servers map[uint][]domain.ServeServer) {
	servers = make(map[uint][]domain.ServeServer)
	res, _ := s.ServeRepo.GetServers(serveIds)
	for _, item := range res {
		var server domain.ServeServer
		copier.CopyWithOption(&server, &item, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		servers[server.ServeId] = append(servers[server.ServeId], server)
	}
	return
}

func (s *DocumentService) GetSecurities(serveIds []uint) (securities map[uint][]domain.ServeSecurityReq) {
	securities = make(map[uint][]domain.ServeSecurityReq)
	res, _ := s.ServeRepo.GetSecurities(serveIds)
	for _, item := range res {
		var security domain.ServeSecurityReq
		copier.CopyWithOption(&security, &item, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		securities[uint(security.ServeId)] = append(securities[uint(security.ServeId)], security)
	}
	return
}

func (s *DocumentService) GetGlobalVars(projectId uint) (globalVars []domain.EnvironmentParam) {
	res, _ := s.EnvironmentRepo.ListGlobalVar(projectId)
	copier.CopyWithOption(&globalVars, &res, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	return
}

func (s *DocumentService) GetDocumentVersionList(projectId uint, needLatest bool) (documents []model.EndpointDocument, err error) {
	if needLatest {
		latestDocument := model.EndpointDocument{
			Name:    "实时版本",
			Version: "latest",
		}
		documents = append(documents, latestDocument)
	}

	documentsTmp, err := s.EndpointDocumentRepo.ListByProject(projectId)
	if err != nil {
		return
	}

	documents = append(documents, documentsTmp...)
	return
}

func (s *DocumentService) Publish(req domain.DocumentVersionReq, projectId uint) (err error) {
	err = s.EndpointSnapshotRepo.BatchCreateSnapshot(req, projectId)
	return
}

func (s *DocumentService) RemoveSnapshot(snapshotId uint) (err error) {
	err = s.EndpointSnapshotRepo.DeleteById(snapshotId)
	return
}

func (s *DocumentService) UpdateSnapshotContent(id uint, endpoint model.Endpoint) (err error) {
	err = s.EndpointSnapshotRepo.UpdateContent(id, endpoint)
	return
}

func (s *DocumentService) UpdateDocument(req domain.UpdateDocumentVersionReq) (err error) {
	err = s.EndpointDocumentRepo.Update(req)
	return
}

func (s *DocumentService) GenerateShareLink(req domain.DocumentShareReq) (link string, err error) {
	encryptValue := strconv.Itoa(int(req.ProjectId)) + "-" + strconv.Itoa(int(req.DocumentId)) + "-" + strconv.Itoa(int(req.EndpointId))
	res, err := commUtils.AesCBCEncrypt([]byte(encryptValue), []byte(EncryptKey))
	link = base64.RawURLEncoding.EncodeToString(res)
	return
}

func (s *DocumentService) DecryptShareLink(link string) (req domain.DocumentShareReq, err error) {
	linkByte, err := base64.RawURLEncoding.DecodeString(link)
	if err != nil {
		return
	}

	decryptValue, err := commUtils.AesCBCDecrypt(linkByte, []byte(EncryptKey))
	if err != nil {
		return
	}

	DocumentShareArr := strings.Split(string(decryptValue), "-")

	projectId, _ := strconv.Atoi(DocumentShareArr[0])
	documentId, _ := strconv.Atoi(DocumentShareArr[1])
	endpointId, _ := strconv.Atoi(DocumentShareArr[2])
	req.ProjectId = uint(projectId)
	req.DocumentId = uint(documentId)
	req.EndpointId = uint(endpointId)

	return
}

func (s *DocumentService) GetEndpointsByShare(projectId, endpointId *uint, serveIds *[]uint, documentId uint) (res map[uint][]domain.EndpointReq, err error) {
	var endpoints []*model.Endpoint
	if documentId != 0 {
		if *endpointId != 0 {
			endpoints, err = s.EndpointSnapshotRepo.GetByDocumentIdAndEndpointId(documentId, *endpointId)
		} else {
			endpoints, err = s.EndpointSnapshotRepo.GetByDocumentId(documentId)
		}
	} else if *projectId != 0 {
		if *endpointId != 0 {
			endpoints, err = s.EndpointRepo.GetByEndpointIds([]uint{*endpointId}, false)
		} else {
			endpoints, err = s.EndpointRepo.GetByProjectId(*projectId, false)
		}
	}
	if err != nil {
		return
	}

	if err != nil {
		return
	}

	res = s.GetEndpointsInfo(projectId, serveIds, endpoints)

	return
}

func (s *DocumentService) ContentByShare(link string) (res domain.DocumentRep, err error) {
	var projectId, documentId, endpointId uint
	var serveIds []uint

	req, err := s.DecryptShareLink(link)
	if err != nil {
		return
	}

	projectId, endpointId, documentId = req.ProjectId, req.EndpointId, req.DocumentId

	endpoints, err := s.GetEndpointsByShare(&projectId, &endpointId, &serveIds, documentId)
	if err != nil {
		return
	}

	var version string
	if documentId == 0 {
		version = "latest"
	} else {
		document, err := s.EndpointDocumentRepo.GetById(documentId)
		if err != nil {
			return res, err
		}
		version = document.Version
		documentId = document.ID
	}

	res = s.GetProject(projectId)

	res.Serves = s.GetServes(serveIds, endpoints)
	res.Version = version
	res.DocumentId = documentId

	return
}

func (s *DocumentService) GetDocumentDetail(documentId, endpointId, interfaceId uint) (res map[string]interface{}, err error) {
	var interfaceDetail model.EndpointInterface

	if documentId == 0 {
		interfaceDetail, err = s.EndpointInterfaceRepo.GetDetail(interfaceId)
	} else {
		interfaceDetail, err = s.EndpointSnapshotRepo.GetInterfaceDetail(documentId, endpointId, interfaceId)
	}

	if err != nil {
		return
	}

	endpoint, err := s.EndpointRepo.Get(interfaceDetail.EndpointId)
	if err != nil {
		return
	}

	serveId := endpoint.ServeId
	serves, err := s.ServeRepo.GetServers([]uint{serveId})
	if err != nil {
		return
	}

	s.EndpointService.SchemaConv(&interfaceDetail, serveId)
	res = make(map[string]interface{})
	res["interface"] = interfaceDetail
	res["servers"] = serves

	return
}
