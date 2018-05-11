package com.ngtesting.platform.vo;

public class IssueCustomFieldOptionVo extends BaseVo {

	private static final long serialVersionUID = 8057353932992599921L;
	private String value;
	private String label;
	private String descr;
	private Integer ordr;
	private Long fieldId;

	public Long getFieldId() {
		return fieldId;
	}

	public void setFieldId(Long fieldId) {
		this.fieldId = fieldId;
	}

	public String getValue() {
        return value;
    }

    public void setValue(String value) {
        this.value = value;
    }

    public String getLabel() {
		return label;
	}

	public void setLabel(String label) {
		this.label = label;
	}

	public String getDescr() {
		return descr;
	}

	public void setDescr(String descr) {
		this.descr = descr;
	}

    public Integer getOrdr() {
        return ordr;
    }

    public void setOrdr(Integer ordr) {
        this.ordr = ordr;
    }
}
