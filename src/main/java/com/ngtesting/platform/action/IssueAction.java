package com.ngtesting.platform.action;

import com.alibaba.fastjson.JSONObject;
import com.ngtesting.platform.config.Constant;
import com.ngtesting.platform.entity.TestCase;
import com.ngtesting.platform.entity.TestCasePriority;
import com.ngtesting.platform.entity.TestCaseType;
import com.ngtesting.platform.service.*;
import com.ngtesting.platform.util.AuthPassport;
import com.ngtesting.platform.vo.*;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.ResponseBody;

import javax.servlet.http.HttpServletRequest;
import java.util.HashMap;
import java.util.List;
import java.util.Map;


@Controller
@RequestMapping(Constant.API_PATH_CLIENT + "issue/")
public class IssueAction extends BaseAction {
	@Autowired
	ProjectService projectService;
	@Autowired
	IssueService issueService;
    @Autowired
    CaseTypeService caseTypeService;
    @Autowired
    CasePriorityService casePriorityService;
	@Autowired
	CustomFieldService customFieldService;

	@AuthPassport(validate = true)
	@RequestMapping(value = "query", method = RequestMethod.POST)
	@ResponseBody
	public Map<String, Object> query(HttpServletRequest request, @RequestBody JSONObject json) {
		Map<String, Object> ret = new HashMap<String, Object>();

        Long orgId = json.getLong("orgId");
		Long projectId = json.getLong("projectId");

		List<TestCase> ls = issueService.query(projectId);
        List<TestCaseVo> vos = issueService.genVos(ls, false);

        List<TestCaseType> caseTypePos = caseTypeService.list(orgId);
        List<CaseTypeVo> caseTypeList = caseTypeService.genVos(caseTypePos);

        List<TestCasePriority> casePriorityPos = casePriorityService.list(orgId);
        List<CasePriorityVo> casePriorityList = casePriorityService.genVos(casePriorityPos);

        List<CustomFieldVo> customFieldList = customFieldService.listForCaseByProject(orgId, projectId);

        ret.put("data", vos);
        ret.put("caseTypeList", caseTypeList);
        ret.put("casePriorityList", casePriorityList);
		ret.put("customFields", customFieldList);
		ret.put("code", Constant.RespCode.SUCCESS.getCode());
		return ret;
	}

	@AuthPassport(validate = true)
	@RequestMapping(value = "queryForSuiteSelection", method = RequestMethod.POST)
	@ResponseBody
	public Map<String, Object> queryForSuiteSelection(HttpServletRequest request, @RequestBody JSONObject json) {
		Map<String, Object> ret = new HashMap<String, Object>();

		Long projectId = json.getLong("projectId");
        Long caseProjectId = json.getLong("caseProjectId");
		Long suiteId = json.getLong("suiteId");

        List<TestCaseVo> vos = issueService.queryForSuiteSelection(projectId, caseProjectId, suiteId);
		List<TestProjectVo> projects = projectService.listBrothers(projectId);

		ret.put("data", vos);
		ret.put("brotherProjects", projects);
		ret.put("code", Constant.RespCode.SUCCESS.getCode());
		return ret;
	}

	@AuthPassport(validate = true)
	@RequestMapping(value = "queryForRunSelection", method = RequestMethod.POST)
	@ResponseBody
	public Map<String, Object> queryForRunSelection(HttpServletRequest request, @RequestBody JSONObject json) {
		Map<String, Object> ret = new HashMap<String, Object>();

		Long projectId = json.getLong("projectId");
        Long caseProjectId = json.getLong("caseProjectId");
		Long runId = json.getLong("runId");

		List<TestCaseVo> vos = issueService.queryForRunSelection(projectId, caseProjectId, runId);
		List<TestProjectVo> projects = projectService.listBrothers(projectId);

		ret.put("data", vos);
		ret.put("brotherProjects", projects);
		ret.put("code", Constant.RespCode.SUCCESS.getCode());
		return ret;
	}

    @AuthPassport(validate = true)
    @RequestMapping(value = "get", method = RequestMethod.POST)
    @ResponseBody
    public Map<String, Object> get(HttpServletRequest request, @RequestBody JSONObject json) {
        Map<String, Object> ret = new HashMap<String, Object>();

        UserVo userVo = (UserVo) request.getSession().getAttribute(Constant.HTTP_SESSION_USER_KEY);
        Long orgId = userVo.getDefaultOrgId();
        Long caseId = json.getLong("id");

        TestCaseVo vo = issueService.getById(caseId);

        ret.put("data", vo);
        ret.put("code", Constant.RespCode.SUCCESS.getCode());
        return ret;
    }

	@AuthPassport(validate = true)
	@RequestMapping(value = "save", method = RequestMethod.POST)
	@ResponseBody
	public Map<String, Object> save(HttpServletRequest request, @RequestBody JSONObject json) {
		Map<String, Object> ret = new HashMap<String, Object>();

        UserVo userVo = (UserVo) request.getSession().getAttribute(Constant.HTTP_SESSION_USER_KEY);

		TestCase po = issueService.save(json, userVo);
		TestCaseVo caseVo = issueService.genVo(po, true);

		ret.put("data", caseVo);
		ret.put("code", Constant.RespCode.SUCCESS.getCode());
		return ret;
	}

	@AuthPassport(validate = true)
	@RequestMapping(value = "rename", method = RequestMethod.POST)
	@ResponseBody
	public Map<String, Object> rename(HttpServletRequest request, @RequestBody JSONObject json) {
		Map<String, Object> ret = new HashMap<String, Object>();

		UserVo userVo = (UserVo) request.getSession().getAttribute(Constant.HTTP_SESSION_USER_KEY);

		TestCase testCasePo = issueService.renamePers(json, userVo);
        issueService.updateParentIfNeededPers(testCasePo.getpId());
		TestCaseVo caseVo = issueService.genVo(testCasePo);

		ret.put("data", caseVo);
		ret.put("code", Constant.RespCode.SUCCESS.getCode());
		return ret;
	}

	@AuthPassport(validate = true)
	@RequestMapping(value = "delete", method = RequestMethod.POST)
	@ResponseBody
	public Map<String, Object> delete(HttpServletRequest request, @RequestBody JSONObject json) {
		Map<String, Object> ret = new HashMap<String, Object>();

		Long id = json.getLong("id");

		UserVo userVo = (UserVo) request.getSession().getAttribute(Constant.HTTP_SESSION_USER_KEY);

		TestCase testCase = issueService.delete(id, userVo);
		issueService.updateParentIfNeededPers(testCase.getpId());

		ret.put("code", Constant.RespCode.SUCCESS.getCode());
		return ret;
	}

	@AuthPassport(validate = true)
	@RequestMapping(value = "move", method = RequestMethod.POST)
	@ResponseBody
	public Map<String, Object> move(HttpServletRequest request, @RequestBody JSONObject json) {
		Map<String, Object> ret = new HashMap<String, Object>();

		UserVo userVo = (UserVo) request.getSession().getAttribute(Constant.HTTP_SESSION_USER_KEY);

        Long srcId = json.getLong("srcId");
        Long parentId = issueService.getById(srcId).getpId();
        Long targetId = json.getLong("targetId");
        TestCaseVo vo = issueService.movePers(json, userVo);

        issueService.updateParentIfNeededPers(parentId);
        issueService.updateParentIfNeededPers(targetId);

		ret.put("data", vo);
		ret.put("code", Constant.RespCode.SUCCESS.getCode());
		return ret;
	}

	@AuthPassport(validate = true)
	@RequestMapping(value = "saveField", method = RequestMethod.POST)
	@ResponseBody
	public Map<String, Object> saveField(HttpServletRequest request, @RequestBody JSONObject json) {
		Map<String, Object> ret = new HashMap<String, Object>();

		UserVo userVo = (UserVo) request.getSession().getAttribute(Constant.HTTP_SESSION_USER_KEY);

		TestCase po = issueService.saveField(json, userVo);
        TestCaseVo caseVo = issueService.genVo(po);

		ret.put("data", caseVo);
		ret.put("code", Constant.RespCode.SUCCESS.getCode());
		return ret;
	}

	@AuthPassport(validate = true)
	@RequestMapping(value = "changeContentType", method = RequestMethod.POST)
	@ResponseBody
	public Map<String, Object> changeContentType(HttpServletRequest request, @RequestBody JSONObject json) {
		Map<String, Object> ret = new HashMap<String, Object>();

		Long id = json.getLong("id");
        String contentType = json.getString("contentType");

		TestCase po = issueService.changeContentTypePers(id, contentType);

		ret.put("code", Constant.RespCode.SUCCESS.getCode());
		return ret;
	}

	@AuthPassport(validate = true)
	@RequestMapping(value = "reviewPass", method = RequestMethod.POST)
	@ResponseBody
	public Map<String, Object> reviewPass(HttpServletRequest request, @RequestBody JSONObject json) {
		Map<String, Object> ret = new HashMap<String, Object>();

		Long id = json.getLong("id");
		Boolean pass = json.getBoolean("pass");

		TestCase po = issueService.reviewPassPers(id, pass);
        TestCaseVo caseVo = issueService.genVo(po);

        ret.put("reviewResult", caseVo);
		ret.put("code", Constant.RespCode.SUCCESS.getCode());
		return ret;
	}

}
