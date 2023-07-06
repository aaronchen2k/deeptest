package service

import (
	v1 "github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/repo"
	"time"
)

type SummaryProjectUserRankingService struct {
	SummaryProjectUserRankingRepo *repo.SummaryProjectUserRankingRepo `inject:""`
}

func (s *SummaryProjectUserRankingService) ProjectUserRanking(cycle int64, projectId int64) (resRankingList v1.ResRankingList, err error) {

	//获取即时数据
	newRankings, _ := s.GetRanking(projectId)

	//查询所有用户名字，map的key为userId
	userInfo, _ := s.FindAllUserName()
	var lastWeekRanking map[int64]model.SummaryProjectUserRanking
	lastWeekRanking = make(map[int64]model.SummaryProjectUserRanking)

	for _, newRanking := range newRankings {
		var resRanking v1.ResUserRanking
		var testTotal int64
		var scenarioTotal int64

		//查询7天前直到现在，最靠前的数据
		earlierDateStartTime, todayEndTime := GetEarlierDateUntilTodayStartAndEndTime(-7)
		lastWeekRanking, err = s.FindMinDataByDateAndProjectIdOfMap(earlierDateStartTime, todayEndTime, projectId)

		if cycle == 1 {
			//全部范围数据,就是最新的数据，newRanking
			scenarioTotal = newRanking.ScenarioTotal
			testTotal = newRanking.TestCaseTotal
		} else if cycle == 0 {
			//当月范围数据
			//先查询31天前的数据
			earlierDateStartTime, earlierDateEndTime := GetEarlierDateStartAndEndTime(-31)
			lastMonthLastDayRanking, _ := s.FindMaxDataByDateAndProjectIdOfMap(earlierDateStartTime, earlierDateEndTime, projectId)
			//那newRanking的所有数据，减去30天前的，就是当月增量数据情况
			scenarioTotal = newRanking.ScenarioTotal - lastMonthLastDayRanking[newRanking.UserId].ScenarioTotal
			testTotal = newRanking.TestCaseTotal - lastMonthLastDayRanking[newRanking.UserId].TestCaseTotal
		}
		if lastWeekRanking[newRanking.UserId].Sort != 0 {
			resRanking.Hb = lastWeekRanking[newRanking.UserId].Sort - newRanking.Sort
		}
		resRanking.Sort = newRanking.Sort
		resRanking.ScenarioTotal = scenarioTotal
		resRanking.TestCaseTotal = testTotal
		resRanking.UserName = userInfo[newRanking.UserId]
		lastUpdateTime, _ := s.FindUserLastUpdateTestCasesByProjectId(projectId)
		if lastUpdateTime[newRanking.UserId] != nil {
			resRanking.UpdatedAt = lastUpdateTime[newRanking.UserId].Format("2006-01-02 15:04:05")
		} else {
			resRanking.UpdatedAt = "------"
		}
		resRanking.UserId = newRanking.UserId
		resRankingList.UserRankingList = append(resRankingList.UserRankingList, resRanking)
	}

	//由于存在当月选项，当月数据需要重新进行排序，不累积
	if len(resRankingList.UserRankingList) != 0 {
		resRankingList, _ = s.SortRankingList(resRankingList)
	}

	return
}

func (s *SummaryProjectUserRankingService) HandlerSummaryProjectUserRankingRepo() *repo.SummaryProjectUserRankingRepo {
	return repo.NewSummaryProjectUserRankingRepo()
}

func (s *SummaryProjectUserRankingService) Create(req model.SummaryProjectUserRanking) (err error) {
	return s.HandlerSummaryProjectUserRankingRepo().Create(req)
}

func (s *SummaryProjectUserRankingService) CreateByDate(req model.SummaryProjectUserRanking) (err error) {
	startTime, endTime := GetTodayStartAndEndTime()
	id, err := s.Existed(startTime, endTime, req.ProjectId, req.UserId)
	if id == 0 {
		err = s.Create(req)
	} else {
		err = s.UpdateColumnsByDate(id, req)
	}
	return
}

func (s *SummaryProjectUserRankingService) UpdateColumnsByDate(id int64, req model.SummaryProjectUserRanking) (err error) {
	return s.HandlerSummaryProjectUserRankingRepo().UpdateColumnsByDate(id, req)
}

func (s *SummaryProjectUserRankingService) FindProjectIds() (projectIds []int64, err error) {
	return s.HandlerSummaryProjectUserRankingRepo().FindProjectIds()
}

func (s *SummaryProjectUserRankingService) Existed(startTime string, endTiem string, projectId int64, userId int64) (id int64, err error) {
	return s.HandlerSummaryProjectUserRankingRepo().Existed(startTime, endTiem, projectId, userId)
}

func (s *SummaryProjectUserRankingService) FindByProjectId(projectId int64) (summaryProjectUserRanking []model.SummaryProjectUserRanking, err error) {
	return s.HandlerSummaryProjectUserRankingRepo().FindByProjectId(projectId)
}

func (s *SummaryProjectUserRankingService) FindMaxDataByDateAndProjectId(startTime string, endTime string, projectId int64) (summaryProjectUserRanking []model.SummaryProjectUserRanking, err error) {

	return s.HandlerSummaryProjectUserRankingRepo().FindMaxDataByDateAndProjectId(startTime, endTime, projectId)
}

func (s *SummaryProjectUserRankingService) FindMinDataByDateAndProjectId(startTime string, endTime string, projectId int64) (summaryProjectUserRanking []model.SummaryProjectUserRanking, err error) {

	return s.HandlerSummaryProjectUserRankingRepo().FindMinDataByDateAndProjectId(startTime, endTime, projectId)
}

func (s *SummaryProjectUserRankingService) FindMaxDataByDateAndProjectIdOfMap(startTime string, endTime string, projectId int64) (result map[int64]model.SummaryProjectUserRanking, err error) {
	summaryProjectUserRanking, _ := s.FindMaxDataByDateAndProjectId(startTime, endTime, projectId)

	result = make(map[int64]model.SummaryProjectUserRanking, len(summaryProjectUserRanking))
	for _, ranking := range summaryProjectUserRanking {
		result[ranking.UserId] = ranking
	}
	return
}

func (s *SummaryProjectUserRankingService) FindMinDataByDateAndProjectIdOfMap(startTime string, endTime string, projectId int64) (result map[int64]model.SummaryProjectUserRanking, err error) {
	summaryProjectUserRanking, _ := s.FindMinDataByDateAndProjectId(startTime, endTime, projectId)

	result = make(map[int64]model.SummaryProjectUserRanking, len(summaryProjectUserRanking))
	for _, ranking := range summaryProjectUserRanking {
		result[ranking.UserId] = ranking
	}
	return
}

func (s *SummaryProjectUserRankingService) CheckUpdated(lastUpdateTime *time.Time) (result bool, err error) {
	return s.HandlerSummaryProjectUserRankingRepo().CheckUpdated(lastUpdateTime)
}

func (s *SummaryProjectUserRankingService) FindScenarioTotalOfUserGroupByProject() (ScenariosTotal map[int64][]model.ProjectUserTotal, err error) {

	results, err := s.HandlerSummaryProjectUserRankingRepo().FindProjectUserScenarioTotal()
	ScenariosTotal = make(map[int64][]model.ProjectUserTotal, len(results))

	for _, result := range results {
		ScenariosTotal[result.ProjectId] = append(ScenariosTotal[result.ProjectId], result)
	}
	return
}

func (s *SummaryProjectUserRankingService) FindTestCasesTotalOfUserGroupByProject() (testCasesTotal map[int64][]model.ProjectUserTotal, err error) {

	results, err := s.HandlerSummaryProjectUserRankingRepo().FindProjectUserTestCasesTotal()
	testCasesTotal = make(map[int64][]model.ProjectUserTotal, len(results))

	for _, result := range results {
		testCasesTotal[result.ProjectId] = append(testCasesTotal[result.ProjectId], result)
	}
	return
}

func (s *SummaryProjectUserRankingService) FindCasesTotalByProjectId(projectId int64) (result map[int64]int64, err error) {

	counts, err := s.HandlerSummaryProjectUserRankingRepo().FindCasesTotalByProjectId(projectId)
	result = make(map[int64]int64, len(counts))
	for _, tmp := range counts {
		result[tmp.CreateUserId] = tmp.Count
	}

	return
}

func (s *SummaryProjectUserRankingService) FindScenariosTotalByProjectId(projectId int64) (result map[int64]int64, err error) {

	counts, err := s.HandlerSummaryProjectUserRankingRepo().FindScenariosTotalByProjectId(projectId)
	result = make(map[int64]int64, len(counts))
	for _, tmp := range counts {
		result[tmp.CreateUserId] = tmp.Count
	}

	return
}

func (s *SummaryProjectUserRankingService) FindUserLastUpdateTestCasesByProjectId(projectId int64) (result map[int64]*time.Time, err error) {

	updateTime, err := s.HandlerSummaryProjectUserRankingRepo().FindUserLastUpdateTestCasesByProjectId(projectId)
	result = make(map[int64]*time.Time, len(updateTime))
	for _, tmp := range updateTime {
		result[tmp.CreatedBy] = tmp.UpdatedAt
	}
	return
}

func (s *SummaryProjectUserRankingService) FindAllUserName() (result map[int64]string, err error) {
	users, err := s.HandlerSummaryProjectUserRankingRepo().FindAllUserName()
	result = make(map[int64]string, len(users))
	for _, user := range users {
		result[user.Id] = user.Name
	}
	return
}

func (s *SummaryProjectUserRankingService) FindUserByProjectId(projectId int64) (users []model.RankingUser, err error) {
	users, err = s.HandlerSummaryProjectUserRankingRepo().FindUserByProjectId(projectId)
	return
}

func (s *SummaryProjectUserRankingService) FindUserIdsByProjectId(projectId int64) (userIds []int64, err error) {
	userIds, err = s.HandlerSummaryProjectUserRankingRepo().FindUserIdsByProjectId(projectId)
	return
}

func (s *SummaryProjectUserRankingService) ForMap(userTotal []model.UserTotal) (ret []map[int64]int64, err error) {

	user := make(map[int64]int64, len(userTotal))
	for _, u := range userTotal {
		user[u.CreateUserId] = u.Count
		ret = append(ret, user)
	}
	return
}

func (s *SummaryProjectUserRankingService) SortRanking(data []model.SummaryProjectUserRanking) (ret []model.SummaryProjectUserRanking, err error) {
	length := len(data)
	for i := 0; i < length; i++ {
		max := data[i].TestCaseTotal

		for x := i + 1; x < length; x++ {
			if data[x].TestCaseTotal > max {

				tmp := data[i]
				data[i] = data[x]
				data[x] = tmp
				data[i].Sort = int64(x)
			}
		}
	}
	data[length-1].Sort = int64(length)
	ret = data
	return
}

func (s *SummaryProjectUserRankingService) SortRankingList(data v1.ResRankingList) (ret v1.ResRankingList, err error) {
	list := data.UserRankingList
	length := len(list)
	for i := 0; i < length; i++ {
		max := list[i].TestCaseTotal
		for x := i + 1; x < length; x++ {
			if list[x].TestCaseTotal > max {
				tmp := list[i]
				list[i] = list[x]
				list[x] = tmp
				list[i].Sort = int64(x)
			}
		}
	}
	list[length-1].Sort = int64(length)
	ret = data
	return
}

func (s *SummaryProjectUserRankingService) GetRanking(projectId int64) (rankings []model.SummaryProjectUserRanking, err error) {

	users, err := s.FindUserIdsByProjectId(projectId)

	cases, err := s.FindCasesTotalByProjectId(projectId)
	scenarios, err := s.FindScenariosTotalByProjectId(projectId)
	lastUpdateTime, _ := s.FindUserLastUpdateTestCasesByProjectId(projectId)

	for _, user := range users {
		var ranking model.SummaryProjectUserRanking
		ranking.UserId = user
		ranking.ProjectId = projectId
		ranking.ScenarioTotal = scenarios[user]
		ranking.TestCaseTotal = cases[user]
		ranking.UpdatedAt = lastUpdateTime[user]
		rankings = append(rankings, ranking)
	}
	if len(rankings) != 0 {
		rankings, err = s.SortRanking(rankings)
	}

	return
}

func (s *SummaryProjectUserRankingService) SaveRanking() (err error) {
	projectIds, err := s.FindProjectIds()
	for _, projectId := range projectIds {
		rankings, _ := s.GetRanking(projectId)
		for _, ranking := range rankings {
			err := s.CreateByDate(ranking)
			if err != nil {
				return err
			}
		}
	}
	return
}
