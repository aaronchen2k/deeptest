package cn.mobiu.events.action.client;

import java.util.HashMap;
import java.util.List;
import java.util.Map;

import javax.servlet.http.HttpServletRequest;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.ResponseBody;

import com.alibaba.fastjson.JSONObject;

import cn.mobiu.events.constants.Constant;
import cn.mobiu.events.entity.EvtBizcard;
import cn.mobiu.events.entity.EvtClient;
import cn.mobiu.events.service.BizcardService;
import cn.mobiu.events.service.DocumentService;
import cn.mobiu.events.util.AuthPassport;
import cn.mobiu.events.util.BeanUtilEx;
import cn.mobiu.events.vo.BizcardVo;
import cn.mobiu.events.vo.DocumentVo;
import cn.mobiu.events.vo.GuestVo;
import cn.mobiu.events.vo.Page;


@Controller
@RequestMapping(Constant.API_PATH_CLIENT + "document/")
public class DocumentAction extends BaseAction {
	@Autowired
	DocumentService documentService;
	
	@AuthPassport(validate = true)
	@RequestMapping(value = "list", method = RequestMethod.POST)
	@ResponseBody
	public Map<String, Object> list(HttpServletRequest request) {
		Map<String, Object> ret = new HashMap<String, Object>();
		JSONObject req = reqJson(request);
		String eventId = req.getString("eventId");
		
		EvtClient client = (EvtClient) request.getSession().getAttribute(Constant.HTTP_SESSION_CLIENT_KEY);

		ret.put("code", Constant.RespCode.SUCCESS.getCode());
		return ret;
	}

	@AuthPassport(validate = true)
	@RequestMapping(value = "listByPage", method = RequestMethod.POST)
	@ResponseBody
	public Map<String, Object> listByPage(HttpServletRequest request, @RequestBody JSONObject json) {
		Map<String, Object> ret = new HashMap<String, Object>();
		
		long eventId = json.getLong("eventId");
		int currentPage = json.getInteger("currentPage") == null? 0: json.getInteger("currentPage") - 1;
		int itemsPerPage = json.getInteger("itemsPerPage") == null? Constant.PAGE_SIZE: json.getInteger("itemsPerPage");
		
		EvtClient client = (EvtClient) request.getSession().getAttribute(Constant.HTTP_SESSION_CLIENT_KEY);
		
		Page page = documentService.listByPage(eventId, currentPage, itemsPerPage, null);
		List<DocumentVo> vos = documentService.genVos(page.getItems());
        
		ret.put("totalItems", page.getTotal());
        ret.put("data", vos);
		ret.put("code", Constant.RespCode.SUCCESS.getCode());
		return ret;
	}
}
