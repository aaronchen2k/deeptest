package com.ngtesting.platform.service;

import com.alibaba.fastjson.JSONObject;
import com.ngtesting.platform.entity.TestPlan;
import com.ngtesting.platform.vo.Page;
import com.ngtesting.platform.vo.TestPlanVo;
import com.ngtesting.platform.vo.UserVo;

import java.util.List;

public interface PlanService extends BaseService {

	Page page(Long projectId, String status, String keywords, Integer currentPage, Integer itemsPerPage);
	TestPlanVo getById(Long caseId);
	TestPlan save(JSONObject json, UserVo optUser);
	TestPlan delete(Long vo, Long userId);

	List<TestPlan> listByOrg(Long orgId);

	List<TestPlan> listByProject(Long projectId, String type);

	List<TestPlanVo> genVos(List<TestPlan> pos);
	TestPlanVo genVo(TestPlan po);

    TestPlan updatePo(TestPlanVo vo);
}
