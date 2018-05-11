package com.ngtesting.platform.service.impl;

import com.ngtesting.platform.entity.TestCaseExeStatus;
import com.ngtesting.platform.service.IssueStatusService;
import com.ngtesting.platform.util.BeanUtilEx;
import com.ngtesting.platform.vo.CaseExeStatusVo;
import org.hibernate.criterion.DetachedCriteria;
import org.hibernate.criterion.Order;
import org.hibernate.criterion.Restrictions;
import org.springframework.stereotype.Service;

import java.util.LinkedList;
import java.util.List;
import java.util.UUID;

@Service
public class IssueStatusServiceImpl extends BaseServiceImpl implements IssueStatusService {
	@Override
	public List<TestCaseExeStatus> list(Long orgId) {
        DetachedCriteria dc = DetachedCriteria.forClass(TestCaseExeStatus.class);

        dc.add(Restrictions.eq("orgId", orgId));
        dc.add(Restrictions.eq("disabled", Boolean.FALSE));
        dc.add(Restrictions.eq("deleted", Boolean.FALSE));

        dc.addOrder(Order.asc("displayOrder"));
        List ls = findAllByCriteria(dc);

		return ls;
	}
	@Override
	public List<CaseExeStatusVo> listVos(Long orgId) {
        List ls = list(orgId);

        List<CaseExeStatusVo> vos = genVos(ls);
		return vos;
	}

	@Override
	public TestCaseExeStatus save(CaseExeStatusVo vo, Long orgId) {
		if (vo == null) {
			return null;
		}

		TestCaseExeStatus po;
		if (vo.getId() != null) {
			po = (TestCaseExeStatus) get(TestCaseExeStatus.class, vo.getId());
		} else {
			po = new TestCaseExeStatus();
		}

		BeanUtilEx.copyProperties(po, vo);

		po.setOrgId(orgId);

		if (vo.getId() == null) {
			po.setCode(UUID.randomUUID().toString());

			String hql = "select max(displayOrder) from TestCaseExeStatus";
			Integer maxOrder = (Integer) getByHQL(hql);
	        po.setDisplayOrder(maxOrder + 10);
		}

		saveOrUpdate(po);
		return po;
	}

	@Override
	public boolean delete(Long id) {
		TestCaseExeStatus po = (TestCaseExeStatus) get(TestCaseExeStatus.class, id);
		po.setDeleted(true);
		saveOrUpdate(po);

		return true;
	}

	@Override
	public boolean changeOrderPers(Long id, String act) {
		TestCaseExeStatus type = (TestCaseExeStatus) get(TestCaseExeStatus.class, id);

        String hql = "from TestCaseExeStatus tp where tp.deleted = false and tp.disabled = false ";
        if ("up".equals(act)) {
        	hql += "and tp.displayOrder < ? order by displayOrder desc";
        } else if ("down".equals(act)) {
        	hql += "and tp.displayOrder > ? order by displayOrder asc";
        } else {
        	return false;
        }

        TestCaseExeStatus neighbor = (TestCaseExeStatus) getDao().findFirstByHQL(hql, type.getDisplayOrder());

        Integer order = type.getDisplayOrder();
        type.setDisplayOrder(neighbor.getDisplayOrder());
        neighbor.setDisplayOrder(order);

        saveOrUpdate(type);
        saveOrUpdate(neighbor);

		return true;
	}

//	@Override
//	public void createDefaultBasicDataPers(Long orgId) {
//		DetachedCriteria dc = DetachedCriteria.forClass(TestCaseExeStatus.class);
//		dc.add(Restrictions.eq("isBuildIn", true));
//		dc.add(Restrictions.eq("disabled", Boolean.FALSE));
//		dc.add(Restrictions.eq("deleted", Boolean.FALSE));
//
//		dc.addOrder(Order.asc("displayOrder"));
//		List<TestCaseExeStatus> ls = findAllByCriteria(dc);
//
//		for (TestCaseExeStatus p : ls) {
//			TestCaseExeStatus temp = new TestCaseExeStatus();
//			BeanUtilEx.copyProperties(temp, p);
//			temp.setId(null);
//			temp.setOrgId(orgId);
//			temp.setBuildIn(false);
//			saveOrUpdate(temp);
//		}
//	}

	@Override
	public CaseExeStatusVo genVo(TestCaseExeStatus po) {
		if (po == null) {
			return null;
		}
		CaseExeStatusVo vo = new CaseExeStatusVo();
		BeanUtilEx.copyProperties(vo, po);

		return vo;
	}
	@Override
	public List<CaseExeStatusVo> genVos(List<TestCaseExeStatus> pos) {
        List<CaseExeStatusVo> vos = new LinkedList<CaseExeStatusVo>();

        for (TestCaseExeStatus po: pos) {
        	CaseExeStatusVo vo = genVo(po);
        	vos.add(vo);
        }
		return vos;
	}


}
