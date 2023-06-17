package model

import (
	serverConsts "github.com/aaronchen2k/deeptest/internal/server/consts"
	"github.com/kataras/iris/v12"
)

type TestInterface struct {
	BaseModel

	Title  string                         `json:"title"`
	Desc   string                         `json:"desc"`
	IsLeaf bool                           `json:"isLeaf"`
	Type   serverConsts.TestInterfaceType `json:"type"`

	ParentId  uint `json:"parentId"`
	ServerId  uint `json:"serverId"`
	ServeId   uint `json:"serveId"`
	ProjectId uint `json:"projectId"`
	UseID     uint `json:"useId"`

	Ordr     int              `json:"ordr"`
	Children []*TestInterface `gorm:"-" json:"children"`
	Slots    iris.Map         `gorm:"-" json:"slots"`

	DebugInterfaceId uint            `json:"debugInterfaceId"`
	DebugInterface   *DebugInterface `gorm:"-" json:"debugInterface"`
}

func (TestInterface) TableName() string {
	return "biz_test_interface"
}