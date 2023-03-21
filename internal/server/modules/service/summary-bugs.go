package service

import (
	"errors"
	"fmt"
	v1 "github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/repo"
	"strconv"
)

type SummaryBugsService struct {
	SummaryBugsRepo *repo.SummaryBugsRepo `inject:""`
}

func NewSummaryBugsService() *SummaryBugsService {
	return &SummaryBugsService{}
}

func (s *SummaryBugsService) Bugs(projectId int64) (res v1.ResSummaryBugs, err error) {

	var summaryBugsSeverity []model.SummaryBugsSeverity
	if projectId == 0 {
		res.Total, err = s.Count()
		summaryBugsSeverity, err = s.FindGroupByBugSeverity()
	} else {
		res.Total, err = s.CountByProjectId(projectId)
		summaryBugsSeverity, err = s.FindByProjectIdGroupByBugSeverity(projectId)
	}

	if err == nil {
		for _, result := range summaryBugsSeverity {
			switch result.BugSeverity {
			case "critical":
				res.Critical = DecimalPer(result.Count, res.Total)
				res.Critical, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", res.Critical), 64)
			case "blocker":
				res.Blocker = DecimalPer(result.Count, res.Total)
				res.Blocker, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", res.Blocker), 64)
			case "deadly":
				res.Deadly = DecimalPer(result.Count, res.Total)
				res.Deadly, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", res.Deadly), 64)
			case "major":
				res.Major = DecimalPer(result.Count, res.Total)
				res.Major, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", res.Major), 64)
			case "minor":
				res.Minor = DecimalPer(result.Count, res.Total)
				res.Minor, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", res.Minor), 64)
			default:
				errors.New("Bug严重程度错误,请检查数据")
			}
		}
	}

	return
}

func DecimalPer(number int64, total int64) float64 {
	value := float64(number) / float64(total)
	return value * 100.0
}

// FindByProjectId
func (s *SummaryBugsService) FindByProjectIdGroupByBugSeverity(projectId int64) (summaryBugsSeverity []model.SummaryBugsSeverity, err error) {
	summaryBugsSeverity, err = s.SummaryBugsRepo.FindByProjectIdGroupByBugSeverity(projectId)
	return
}

// FindGroupByBugSeverity
func (s *SummaryBugsService) FindGroupByBugSeverity() (summaryBugsSeverity []model.SummaryBugsSeverity, err error) {
	summaryBugsSeverity, err = s.SummaryBugsRepo.FindGroupByBugSeverity()
	return
}

// Create
func (s *SummaryBugsService) CreateByDate(req model.SummaryBugs) (err error) {
	return s.SummaryBugsRepo.CreateByDate(req)
}

// Count
func (s *SummaryBugsService) Count() (count int64, err error) {
	return s.SummaryBugsRepo.Count()
}

// CountByProjectId
func (s *SummaryBugsService) CountByProjectId(projectId int64) (count int64, err error) {
	return s.SummaryBugsRepo.CountByProjectId(projectId)
}
