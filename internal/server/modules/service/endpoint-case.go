package service

import (
	serverDomain "github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	serverConsts "github.com/aaronchen2k/deeptest/internal/server/consts"
	"github.com/aaronchen2k/deeptest/internal/server/modules/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/repo"
	"github.com/jinzhu/copier"
	"github.com/kataras/iris/v12"
	uuid "github.com/satori/go.uuid"
)

type EndpointCaseService struct {
	EndpointCaseRepo      *repo.EndpointCaseRepo      `inject:""`
	EndpointInterfaceRepo *repo.EndpointInterfaceRepo `inject:""`
	ServeServerRepo       *repo.ServeServerRepo       `inject:""`
	DebugInterfaceRepo    *repo.DebugInterfaceRepo    `inject:""`
	EndpointRepo          *repo.EndpointRepo          `inject:""`
	PreConditionRepo      *repo.PreConditionRepo      `inject:""`
	PostConditionRepo     *repo.PostConditionRepo     `inject:""`
	CategoryRepo          *repo.CategoryRepo          `inject:""`

	EndpointService       *EndpointService       `inject:""`
	DebugInterfaceService *DebugInterfaceService `inject:""`
}

func (s *EndpointCaseService) List(endpointId uint) (ret []model.EndpointCase, err error) {
	ret, err = s.EndpointCaseRepo.List(endpointId)

	return
}

func (s *EndpointCaseService) Get(id int) (ret model.EndpointCase, err error) {
	ret, err = s.EndpointCaseRepo.Get(uint(id))
	// its debug data will load in webpage

	return
}

func (s *EndpointCaseService) Save(req serverDomain.EndpointCaseSaveReq) (casePo model.EndpointCase, err error) {
	s.CopyValueFromRequest(&casePo, req)

	endpoint, err := s.EndpointRepo.Get(req.EndpointId)

	var server model.ServeServer
	if endpoint.ServerId > 0 {
		server, _ = s.ServeServerRepo.Get(endpoint.ServerId)
	} else {
		server, _ = s.ServeServerRepo.GetDefaultByServe(endpoint.ServeId)
	}

	// create new DebugInterface
	url := req.DebugData.Url
	if url == "" {
		url = endpoint.Path
	}

	debugInterface := model.DebugInterface{
		InterfaceBase: model.InterfaceBase{
			Name: req.Name,

			InterfaceConfigBase: model.InterfaceConfigBase{
				Method: consts.GET,
				Url:    url,
			},
		},
		ServeId:  endpoint.ServeId,
		ServerId: server.ID,
		BaseUrl:  server.Url,
	}

	err = s.DebugInterfaceRepo.Save(&debugInterface)

	casePo.ProjectId = endpoint.ProjectId
	casePo.ServeId = endpoint.ServeId
	casePo.DebugInterfaceId = debugInterface.ID
	err = s.EndpointCaseRepo.Save(&casePo)

	if casePo.DebugInterfaceId > 0 {
		values := map[string]interface{}{
			"case_interface_id": casePo.ID,
		}
		err = s.DebugInterfaceRepo.UpdateDebugInfo(casePo.DebugInterfaceId, values)
	}

	return
}

func (s *EndpointCaseService) Copy(id int, userId uint, userName string) (po model.EndpointCase, err error) {
	endpointCase, _ := s.EndpointCaseRepo.Get(uint(id))
	debugData, _ := s.DebugInterfaceService.GetDebugDataFromDebugInterface(endpointCase.DebugInterfaceId)

	req := serverDomain.EndpointCaseSaveReq{
		Name:       "copy-" + endpointCase.Name,
		EndpointId: endpointCase.EndpointId,
		ServeId:    endpointCase.ServeId,
		ProjectId:  endpointCase.ProjectId,

		CreateUserId:   userId,
		CreateUserName: userName,

		DebugData: debugData,
	}

	s.CopyValueFromRequest(&po, req)

	endpoint, err := s.EndpointRepo.Get(req.EndpointId)

	// create new DebugInterface
	url := req.DebugData.Url
	if url == "" {
		url = endpoint.Path
	}

	debugInterface := model.DebugInterface{}

	s.DebugInterfaceService.CopyValueFromRequest(&debugInterface, req.DebugData)
	debugInterface.Name = req.Name
	debugInterface.Url = url

	err = s.DebugInterfaceRepo.Save(&debugInterface)

	// clone conditions
	s.PreConditionRepo.CloneAll(req.DebugData.DebugInterfaceId, 0, debugInterface.ID)
	s.PostConditionRepo.CloneAll(req.DebugData.DebugInterfaceId, 0, debugInterface.ID)

	// save case
	po.ProjectId = endpoint.ProjectId
	po.ServeId = endpoint.ServeId
	po.DebugInterfaceId = debugInterface.ID
	err = s.EndpointCaseRepo.Save(&po)

	if po.DebugInterfaceId > 0 {
		values := map[string]interface{}{
			"case_interface_id": po.ID,
		}
		err = s.DebugInterfaceRepo.UpdateDebugInfo(po.DebugInterfaceId, values)
	}

	return
}

func (s *EndpointCaseService) SaveFromDebugInterface(req serverDomain.EndpointCaseSaveReq) (po model.EndpointCase, err error) {
	debugData := req.DebugData

	// save debug data
	req.DebugData.UsedBy = consts.CaseDebug
	debugInterface, err := s.DebugInterfaceService.SaveAs(debugData)

	// save case
	s.CopyValueFromRequest(&po, req)

	if po.EndpointId == 0 {
		endpointInterface, _ := s.EndpointInterfaceRepo.Get(uint(req.EndpointInterfaceId))
		po.EndpointId = endpointInterface.EndpointId
	}
	endpoint, err := s.EndpointRepo.Get(po.EndpointId)
	po.ProjectId = endpoint.ProjectId
	po.ServeId = endpoint.ServeId

	po.DebugInterfaceId = debugInterface.ID
	po.ID = 0
	err = s.EndpointCaseRepo.Save(&po)

	if po.DebugInterfaceId > 0 {
		values := map[string]interface{}{
			"case_interface_id": po.ID,
		}
		err = s.DebugInterfaceRepo.UpdateDebugInfo(po.DebugInterfaceId, values)
	}

	if err != nil {
		return
	}

	return
}

func (s *EndpointCaseService) UpdateName(req serverDomain.EndpointCaseSaveReq) (err error) {
	err = s.EndpointCaseRepo.UpdateName(req)

	return
}

func (s *EndpointCaseService) Remove(id uint) (err error) {
	err = s.EndpointCaseRepo.Remove(id)
	return
}

func (s *EndpointCaseService) CopyValueFromRequest(po *model.EndpointCase, req serverDomain.EndpointCaseSaveReq) {
	copier.CopyWithOption(po, req, copier.Option{
		DeepCopy: true,
	})
}

func (s *EndpointCaseService) LoadTree(projectId, serveId uint) (ret []*serverDomain.EndpointCaseTree, err error) {
	root, err := s.GetTree(projectId, serveId)

	s.mountCount(root, projectId, serveId)
	if root != nil && len(root.Children) > 0 && root.Children[0] != nil {
		ret = root.Children[0].Children
	}

	return
}

func (s *EndpointCaseService) GetTree(projectId, serveId uint) (root *serverDomain.EndpointCaseTree, err error) {
	categories, err := s.CategoryRepo.ListByProject(serverConsts.EndpointCategory, projectId, 0)
	if err != nil || len(categories) == 0 {
		return
	}
	categoryTos := s.CategoryToTos(categories)

	categoryTos = append(categoryTos, &serverDomain.EndpointCaseTree{Key: -1, Id: uuid.NewV4(), Name: "未分类", ParentId: int64(categories[0].ID), Slots: iris.Map{"icon": "icon"}})

	endpoints, err := s.EndpointRepo.ListByProjectIdAndServeId(projectId, serveId, false)
	if err != nil {
		return
	}
	endpointTos := s.EndpointToTos(endpoints)

	cases, err := s.EndpointCaseRepo.ListByProjectIdAndServeId(projectId, serveId)
	casesTos := s.EndpointCaseToTos(cases)
	if err != nil {
		return
	}

	for _, endpoint := range endpointTos {
		s.makeTree(casesTos, endpoint, serverConsts.EndpointCaseTreeTypeCase)
	}
	for _, category := range categoryTos {
		s.makeTree(endpointTos, category, serverConsts.EndpointCaseTreeTypeEndpoint)
	}
	root = &serverDomain.EndpointCaseTree{}
	s.makeTree(categoryTos, root, serverConsts.EndpointCaseTreeTypeDir)

	return
}

func (s *EndpointCaseService) CategoryToTos(pos []*model.Category) (tos []*serverDomain.EndpointCaseTree) {
	for _, po := range pos {
		to := s.CategoryToTo(po)

		tos = append(tos, to)
	}

	return
}

func (s *EndpointCaseService) CategoryToTo(po *model.Category) (to *serverDomain.EndpointCaseTree) {
	to = &serverDomain.EndpointCaseTree{
		Id:        uuid.NewV4(),
		Key:       int64(po.ID),
		Name:      po.Name,
		Desc:      po.Desc,
		Type:      serverConsts.EndpointCaseTreeTypeDir,
		IsDir:     true,
		ParentId:  int64(po.ParentId),
		ProjectId: po.ProjectId,
		ServeId:   po.ServeId,
	}

	return
}

func (s *EndpointCaseService) EndpointToTos(pos []*model.Endpoint) (tos []*serverDomain.EndpointCaseTree) {
	for _, po := range pos {
		to := s.EndpointToTo(po)

		tos = append(tos, to)
	}

	return
}

func (s *EndpointCaseService) EndpointToTo(po *model.Endpoint) (to *serverDomain.EndpointCaseTree) {
	to = &serverDomain.EndpointCaseTree{
		Id:         uuid.NewV4(),
		Key:        int64(po.ID),
		Name:       po.Title,
		Desc:       po.Description,
		Type:       serverConsts.EndpointCaseTreeTypeEndpoint,
		IsDir:      true,
		CategoryId: uint(po.CategoryId),
		ProjectId:  po.ProjectId,
		ServeId:    po.ServeId,
	}

	return
}

func (s *EndpointCaseService) EndpointCaseToTos(pos []*model.EndpointCase) (tos []*serverDomain.EndpointCaseTree) {
	for _, po := range pos {
		to := s.EndpointCaseToTo(po)

		tos = append(tos, to)
	}

	return
}

func (s *EndpointCaseService) EndpointCaseToTo(po *model.EndpointCase) (to *serverDomain.EndpointCaseTree) {
	to = &serverDomain.EndpointCaseTree{
		Id:               uuid.NewV4(),
		Key:              int64(po.ID),
		Name:             po.Name,
		Desc:             po.Desc,
		Type:             serverConsts.EndpointCaseTreeTypeCase,
		IsDir:            false,
		EndpointId:       po.EndpointId,
		DebugInterfaceId: po.DebugInterfaceId,
		CaseInterfaceId:  po.ID,
		ProjectId:        po.ProjectId,
		ServeId:          po.ServeId,
	}

	return
}

func (s *EndpointCaseService) makeTree(findIn []*serverDomain.EndpointCaseTree, parent *serverDomain.EndpointCaseTree, typ serverConsts.EndpointCaseTreeType) { //参数为父节点，添加父节点的子节点指针切片
	children, _ := s.hasChild(findIn, parent, typ) // 判断节点是否有子节点并返回

	if children != nil {
		parent.Children = append(parent.Children, children[0:]...) // 添加子节点

		for _, child := range children { // 查询子节点的子节点，并添加到子节点
			_, has := s.hasChild(findIn, child, typ)
			if has {
				s.makeTree(findIn, child, typ) // 递归添加节点
			}
		}
	}
}

func (s *EndpointCaseService) hasChild(categories []*serverDomain.EndpointCaseTree, parent *serverDomain.EndpointCaseTree, typ serverConsts.EndpointCaseTreeType) (
	ret []*serverDomain.EndpointCaseTree, yes bool) {

	for _, item := range categories {
		if s.isChild(item, parent, typ) {
			item.Slots = iris.Map{"icon": "icon"}
			//item.Parent = parent // loop json

			ret = append(ret, item)
		}
	}

	if ret != nil {
		yes = true
	}

	return
}

func (s *EndpointCaseService) isChild(child, parent *serverDomain.EndpointCaseTree, typ serverConsts.EndpointCaseTreeType) (res bool) {
	if child == nil || parent == nil {
		return
	}
	switch typ {
	case serverConsts.EndpointCaseTreeTypeDir:
		res = child.ParentId == parent.Key
	case serverConsts.EndpointCaseTreeTypeEndpoint:
		res = int64(child.CategoryId) == parent.Key
	case serverConsts.EndpointCaseTreeTypeCase:
		res = child.EndpointId == uint(parent.Key)
	}

	return
}

func (s *EndpointCaseService) mountCount(root *serverDomain.EndpointCaseTree, projectId, serveId uint) {
	endpointCount, err := s.EndpointCaseRepo.GetEndpointCount(projectId, serveId)
	if err != nil || len(endpointCount) == 0 {
		return
	}

	result := s.convertMap(endpointCount)
	s.mountCountOnNode(root, result)
}
func (s *EndpointCaseService) convertMap(data []serverDomain.EndpointCount) (result map[int64]int64) {
	result = make(map[int64]int64)
	for _, item := range data {
		result[item.EndpointId] = item.Count
	}
	return
}

func (s *EndpointCaseService) mountCountOnNode(root *serverDomain.EndpointCaseTree, data map[int64]int64) int64 {
	switch root.Type {
	case serverConsts.EndpointCaseTreeTypeDir:
		root.Count = 0
	case serverConsts.EndpointCaseTreeTypeEndpoint:
		root.Count = data[root.Key]
	}
	for _, children := range root.Children {
		root.Count += s.mountCountOnNode(children, data)
	}
	return root.Count
}
