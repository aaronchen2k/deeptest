package service

import (
	v1 "github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/server/modules/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/repo"
	dateUtils "github.com/aaronchen2k/deeptest/pkg/lib/date"
	_fileUtils "github.com/aaronchen2k/deeptest/pkg/lib/file"
	logUtils "github.com/aaronchen2k/deeptest/pkg/lib/log"
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/helper/dir"
	"github.com/xuri/excelize/v2"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"
)

type DatapoolService struct {
	DatapoolRepo *repo.DatapoolRepo `inject:""`
}

func NewDatapoolService() *DatapoolService {
	return &DatapoolService{}
}

func (s *DatapoolService) List(projectId uint) (ret []v1.DatapoolReq, err error) {
	ret, err = s.DatapoolRepo.List(projectId)

	return
}

func (s *DatapoolService) Get(id uint) (model.Datapool, error) {
	return s.DatapoolRepo.Get(id)
}

func (s *DatapoolService) Save(req *model.Datapool) (err error) {
	return s.DatapoolRepo.Save(req)
}

func (s *DatapoolService) Delete(id uint) (err error) {
	return s.DatapoolRepo.Delete(id)
}

// Upload 上传文件
func (s *DatapoolService) Upload(ctx iris.Context, fh *multipart.FileHeader) (ret v1.DatapoolUploadResp, err error) {
	filename, err := _fileUtils.GetUploadFileName(fh.Filename)
	if err != nil {
		logUtils.Errorf("获取文件名失败，错误%s", err.Error())
		return
	}

	targetDir := filepath.Join(consts.DirData, dateUtils.DateStr(time.Now()))
	absDir := filepath.Join(dir.GetCurrentAbPath(), targetDir)

	err = dir.InsureDir(targetDir)
	if err != nil {
		logUtils.Errorf("文件上传失败，错误%s", err.Error())
		return
	}

	pth := filepath.Join(absDir, filename)
	_, err = ctx.SaveFormFile(fh, pth)
	if err != nil {
		logUtils.Errorf("文件上传失败，错误%s", "保存文件到本地")
		return
	}

	ret.Path = filepath.Join(targetDir, filename)
	ret.Data, _ = s.ReadExcel(ret.Path)

	return
}

func (s *DatapoolService) ReadExcel(pth string) (ret [][]interface{}, err error) {
	excl, err := excelize.OpenFile(pth)
	if err != nil {
		logUtils.Info("read upload file as excel failed")
		return
	}

	sht := excl.GetSheetList()[0]
	lines, err := excl.GetRows(sht)

	for _, line := range lines {
		var row []interface{}
		for _, col := range line {
			row = append(row, strings.TrimSpace(col))
		}

		ret = append(ret, row)
	}

	return
}
