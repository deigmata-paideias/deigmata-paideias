package indi.yuluo;

import indi.yuluo.service.MySQLDataService;
import jakarta.annotation.Resource;

import org.springframework.ai.tool.annotation.Tool;
import org.springframework.stereotype.Service;

/**
 * @author yuluo
 * @author <a href="mailto:yuluo08290126@gmail.com">yuluo</a>
 */

@Service
public class ToolsDefinitionService {

	@Resource
	private MySQLDataService dataService;

	@Tool(description = "调用此工具函数，将参数转为全大写并返回")
	public String test(String test) {

		return test.toLowerCase();
	}

	@Tool(description = "调用此工具函数，从 mysql 数据查询指定姓名的联系人信息")
	public String queryContact(String name) {

		return dataService.getContextData(name);
	}

}
