package com.ngtesting.platform.service;

import com.alibaba.fastjson.JSONArray;
import com.ngtesting.platform.entity.TestUser;
import com.ngtesting.platform.vo.Page;
import com.ngtesting.platform.vo.RelationOrgGroupUserVo;
import com.ngtesting.platform.vo.UserVo;

import java.util.List;
import java.util.Map;

public interface UserService extends BaseService {
	Page listByPage(Long orgId, String keywords, String disabled, Integer currentPage, Integer itemsPerPage);
	List<Map> getProjectUsers(Long orgId, Long projectId);
    List<TestUser> search(Long orgId, String keywords, JSONArray exceptIds);

	TestUser save(UserVo vo, Long orgId);
	TestUser setLeftSizePers(Long userId, Integer left, String prop);

    List<TestUser> listAllOrgUsers(Long orgId);

    TestUser invitePers(UserVo uerVo, UserVo newUserVo, List<RelationOrgGroupUserVo> relations);
	boolean disable(Long userId, Long orgId);
	boolean removeUserFromOrgPers(Long userId, Long orgId);

	List<UserVo> genVos(List<TestUser> pos);
	UserVo genVo(TestUser user);

}

