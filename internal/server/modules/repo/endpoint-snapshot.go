package repo

import (
	"encoding/json"
	v1 "github.com/deeptest-com/deeptest/cmd/server/v1/domain"
	"github.com/deeptest-com/deeptest/internal/pkg/consts"
	"github.com/deeptest-com/deeptest/internal/server/modules/model"
	logUtils "github.com/deeptest-com/deeptest/pkg/lib/log"
	"go.uber.org/zap"
)

type EndpointSnapshotRepo struct {
	*BaseRepo             `inject:""`
	EndpointInterfaceRepo *EndpointInterfaceRepo `inject:""`
	ServeRepo             *ServeRepo             `inject:""`
	ProjectRepo           *ProjectRepo           `inject:""`
	EndpointRepo          *EndpointRepo          `inject:""`
	EndpointDocumentRepo  *EndpointDocumentRepo  `inject:""`
}

func NewEndpointSnapshotRepo() *EndpointSnapshotRepo {
	return &EndpointSnapshotRepo{}
}

func (r *EndpointSnapshotRepo) BatchCreateSnapshot(tenantId consts.TenantId, req v1.DocumentVersionReq, projectId uint) (documentId uint, err error) {
	documentId, err = r.EndpointDocumentRepo.GetIdByVersionAndProject(tenantId, req, projectId)
	if err != nil {
		return
	}

	if err = r.BatchDeleteByEndpointId(tenantId, req.EndpointIds, documentId); err != nil {
		return
	}

	snapshots := make([]*model.EndpointSnapshot, 0)
	for _, v := range req.EndpointIds {
		endpoint, err := r.EndpointRepo.GetAll(tenantId, v, "v0.1.0")
		if err != nil {
			logUtils.Errorf("create endpoint snapshot error", zap.String("error:", err.Error()), zap.Uint("endpointId:", v))
			continue
		}
		content, _ := json.Marshal(endpoint)

		snapshotTmp := model.EndpointSnapshot{
			EndpointId: endpoint.ID,
			DocumentId: documentId,
			Content:    string(content),
		}
		snapshots = append(snapshots, &snapshotTmp)
	}

	err = r.GetDB(tenantId).Create(snapshots).Error

	return
}

func (r *EndpointSnapshotRepo) GetByDocumentId(tenantId consts.TenantId, documentId uint) (endpoints []*model.Endpoint, err error) {
	var snapshots []model.EndpointSnapshot
	err = r.GetDB(tenantId).Where("document_id = ? and not deleted and not disabled", documentId).Find(&snapshots).Error
	if err != nil {
		return
	}

	for _, v := range snapshots {
		var endpoint model.Endpoint
		_ = json.Unmarshal([]byte(v.Content), &endpoint)
		endpoints = append(endpoints, &endpoint)
	}

	return
}

func (r *EndpointSnapshotRepo) GetByDocumentIdAndEndpointId(tenantId consts.TenantId, documentId, endpointId uint) (endpoints []*model.Endpoint, err error) {
	var snapshot model.EndpointSnapshot
	err = r.GetDB(tenantId).Where("document_id = ?", documentId).
		Where("endpoint_id = ? and not deleted and not disabled", endpointId).
		Find(&snapshot).Error
	if err != nil {
		return
	}

	var endpoint model.Endpoint
	_ = json.Unmarshal([]byte(snapshot.Content), &endpoint)
	endpoints = append(endpoints, &endpoint)

	return
}

func (r *EndpointSnapshotRepo) DeleteById(tenantId consts.TenantId, id uint) (err error) {
	err = r.GetDB(tenantId).
		Where("id = ?", id).
		Delete(&model.EndpointSnapshot{}).Error

	if err != nil {
		logUtils.Errorf("delete endpoint snapshot by id error", zap.String("error:", err.Error()))
		return
	}

	return
}

func (r *EndpointSnapshotRepo) BatchDeleteByEndpointId(tenantId consts.TenantId, endpointIds []uint, documentId uint) (err error) {
	err = r.GetDB(tenantId).
		Where("endpoint_id IN (?)", endpointIds).
		Where("document_id = ?", documentId).
		Delete(&model.EndpointSnapshot{}).Error

	if err != nil {
		logUtils.Errorf("delete endpoint snapshot by id error", zap.String("error:", err.Error()))
		return
	}

	return
}

func (r *EndpointSnapshotRepo) UpdateContent(tenantId consts.TenantId, id uint, endpoint model.Endpoint) (err error) {
	content, _ := json.Marshal(endpoint)
	err = r.GetDB(tenantId).Model(model.EndpointSnapshot{}).
		Where("id = ?", id).
		Update("content = ?", string(content)).Error

	return
}

func (r *EndpointSnapshotRepo) GetContentByDocumentAndEndpoint(tenantId consts.TenantId, documentId, endpointId uint) (endpoint model.Endpoint, err error) {
	var snapshot model.EndpointSnapshot
	err = r.GetDB(tenantId).Where("document_id = ?", documentId).
		Where("endpoint_id = ? and not deleted and not disabled", endpointId).
		Find(&snapshot).Error
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(snapshot.Content), &endpoint)

	return
}

func (r *EndpointSnapshotRepo) GetInterfaceDetail(tenantId consts.TenantId, documentId, endpointId, interfaceId uint) (interf model.EndpointInterface, err error) {
	snapshotContent, err := r.GetContentByDocumentAndEndpoint(tenantId, documentId, endpointId)
	if err != nil {
		return
	}

	for _, v := range snapshotContent.Interfaces {
		if v.ID == interfaceId {
			interf = v
			break
		}
	}
	return
}
