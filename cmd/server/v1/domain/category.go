package serverDomain

import (
	serverConsts "github.com/aaronchen2k/deeptest/internal/server/consts"
	"github.com/kataras/iris/v12"
)

// category
type Category struct {
	Id       int64       `json:"id"`
	Name     string      `json:"name"`
	Desc     string      `json:"desc"`
	ParentId int64       `json:"parentId"`
	Children []*Category `json:"children"`
	Slots    iris.Map    `json:"slots"`
}

type CategoryCreateReq struct {
	Name      string                             `json:"name"`
	Mode      string                             `json:"mode"`
	Type      serverConsts.CategoryDiscriminator `json:"type"`
	ServeId   uint                               `json:"serveId"`
	ModuleId  string                             `json:"moduleId"`
	TargetId  uint                               `json:"targetId"`
	ProjectId uint                               `json:"projectId"`
}

type CategoryReq struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Desc   string `json:"desc"`
	Parent uint
}

type CategoryMoveReq struct {
	Type    serverConsts.CategoryDiscriminator `json:"type"`
	DragKey int                                `json:"dragKey"`
	DropKey int                                `json:"dropKey"`
	DropPos serverConsts.DropPos               `json:"dropPos"`
}
