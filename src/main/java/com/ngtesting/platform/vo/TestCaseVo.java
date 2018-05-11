package com.ngtesting.platform.vo;

import java.util.LinkedList;
import java.util.List;

public class TestCaseVo extends BaseVo {
	private static final long serialVersionUID = -5955583523485410239L;

	private String name;
	private String priority;
	private Integer estimate;
    private String type;
	private String objective;
    private String descr;
    private Integer ordr;
	private Long pId;
	private Boolean isLeaf;
	private Long projectId;

	private Long createById;
	private Long updateById;

    private Boolean checked;

	private String contentType;
	private String content;
	private String key;
	private Boolean reviewResult;

	private String prop01;
	private String prop02;
	private String prop03;
	private String prop04;
	private String prop05;

	private String prop06;
	private String prop07;
	private String prop08;
	private String prop09;
	private String prop10;

	private String prop11;
	private String prop12;
	private String prop13;
	private String prop14;
	private String prop15;

	private String prop16;
	private String prop17;
	private String prop18;
	private String prop19;
	private String prop20;

    private List<TestCaseVo> children = new LinkedList<>();
	private List<TestCaseStepVo> steps = new LinkedList<>();
	private List<TestCaseCommentsVo> comments = new LinkedList<>();
	private List<TestCaseHistoryVo> histories = new LinkedList<>();

	private List<TestCaseAttachmentVo> attachments = new LinkedList<>();

	public List<TestCaseAttachmentVo> getAttachments() {
		return attachments;
	}

	public void setAttachments(List<TestCaseAttachmentVo> attachments) {
		this.attachments = attachments;
	}

	public TestCaseVo() {
        super();
    }

	public List<TestCaseHistoryVo> getHistories() {
		return histories;
	}

	public void setHistories(List<TestCaseHistoryVo> histories) {
		this.histories = histories;
	}

	public Boolean getLeaf() {
		return isLeaf;
	}

	public void setLeaf(Boolean leaf) {
		isLeaf = leaf;
	}

	public List<TestCaseCommentsVo> getComments() {
		return comments;
	}

	public void setComments(List<TestCaseCommentsVo> comments) {
		this.comments = comments;
	}

	public Boolean getReviewResult() {
		return reviewResult;
	}

	public void setReviewResult(Boolean reviewResult) {
		this.reviewResult = reviewResult;
	}

	public String getKey() {
		return key;
	}

	public void setKey(String key) {
		this.key = key;
	}

	public String getContentType() {
		return contentType;
	}

	public void setContentType(String contentType) {
		this.contentType = contentType;
	}

	public String getContent() {
		return content;
	}

	public void setContent(String content) {
		this.content = content;
	}

	public Boolean getChecked() {
        return checked;
    }

    public void setChecked(Boolean checked) {
        this.checked = checked;
    }

    public Long getCreateById() {
        return createById;
    }

    public void setCreateById(Long createById) {
        this.createById = createById;
    }

    public Long getUpdateById() {
        return updateById;
    }

    public void setUpdateById(Long updateById) {
        this.updateById = updateById;
    }

    public List<TestCaseVo> getChildren() {
        return children;
    }

    public void setChildren(List<TestCaseVo> children) {
        this.children = children;
    }

    public Long getProjectId() {
		return projectId;
	}

	public void setProjectId(Long projectId) {
		this.projectId = projectId;
	}

	public Integer getOrdr() {
        return ordr;
    }

    public void setOrdr(Integer ordr) {
        this.ordr = ordr;
    }

    public String getType() {
        return this.type;
    }

    public void setType(String type) {
        this.type = type;
    }

    public String getName() {
		return name;
	}

	public void setName(String name) {
		this.name = name;
	}

	public Long getpId() {
		return pId;
	}

	public void setpId(Long pId) {
		this.pId = pId;
	}

	public String getProp01() {
		return prop01;
	}

	public void setProp01(String prop01) {
		this.prop01 = prop01;
	}

	public String getProp02() {
		return prop02;
	}

	public void setProp02(String prop02) {
		this.prop02 = prop02;
	}

	public String getProp03() {
		return prop03;
	}

	public void setProp03(String prop03) {
		this.prop03 = prop03;
	}

	public String getProp04() {
		return prop04;
	}

	public void setProp04(String prop04) {
		this.prop04 = prop04;
	}

	public String getProp05() {
		return prop05;
	}

	public void setProp05(String prop05) {
		this.prop05 = prop05;
	}

	public String getProp06() {
		return prop06;
	}

	public void setProp06(String prop06) {
		this.prop06 = prop06;
	}

	public String getProp07() {
		return prop07;
	}

	public void setProp07(String prop07) {
		this.prop07 = prop07;
	}

	public String getProp08() {
		return prop08;
	}

	public void setProp08(String prop08) {
		this.prop08 = prop08;
	}

	public String getProp09() {
		return prop09;
	}

	public void setProp09(String prop09) {
		this.prop09 = prop09;
	}

	public String getProp10() {
		return prop10;
	}

	public void setProp10(String prop10) {
		this.prop10 = prop10;
	}

	public String getProp11() {
		return prop11;
	}

	public void setProp11(String prop11) {
		this.prop11 = prop11;
	}

	public String getProp12() {
		return prop12;
	}

	public void setProp12(String prop12) {
		this.prop12 = prop12;
	}

	public String getProp13() {
		return prop13;
	}

	public void setProp13(String prop13) {
		this.prop13 = prop13;
	}

	public String getProp14() {
		return prop14;
	}

	public void setProp14(String prop14) {
		this.prop14 = prop14;
	}

	public String getProp15() {
		return prop15;
	}

	public void setProp15(String prop15) {
		this.prop15 = prop15;
	}

	public String getProp16() {
		return prop16;
	}

	public void setProp16(String prop16) {
		this.prop16 = prop16;
	}

	public String getProp17() {
		return prop17;
	}

	public void setProp17(String prop17) {
		this.prop17 = prop17;
	}

	public String getProp18() {
		return prop18;
	}

	public void setProp18(String prop18) {
		this.prop18 = prop18;
	}

	public String getProp19() {
		return prop19;
	}

	public void setProp19(String prop19) {
		this.prop19 = prop19;
	}

	public String getProp20() {
		return prop20;
	}

	public void setProp20(String prop20) {
		this.prop20 = prop20;
	}

	public String getPriority() {
		return priority;
	}

	public void setPriority(String priority) {
		this.priority = priority;
	}

	public Integer getEstimate() {
		return estimate;
	}

	public void setEstimate(Integer estimate) {
		this.estimate = estimate;
	}

	public String getDescr() {
		return descr;
	}

	public void setDescr(String descr) {
		this.descr = descr;
	}

	public String getObjective() {
		return objective;
	}

	public void setObjective(String objective) {
		this.objective = objective;
	}

//    public List<TestCasePropVo> getProps() {
//        return props;
//    }
//
//    public void setProps(List<TestCasePropVo> props) {
//        this.props = props;
//    }

	public List<TestCaseStepVo> getSteps() {
		return steps;
	}

	public void setSteps(List<TestCaseStepVo> steps) {
		this.steps = steps;
	}
}
