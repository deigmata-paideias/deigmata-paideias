package indi.yuluo;

import org.mybatis.spring.annotation.MapperScan;

import org.springframework.ai.tool.ToolCallbackProvider;
import org.springframework.ai.tool.method.MethodToolCallbackProvider;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.Bean;

/**
 * @author yuluo
 * @author <a href="mailto:yuluo08290126@gmail.com">yuluo</a>
 */

@SpringBootApplication
@MapperScan("indi.yuluo.repository")
public class MCPServerApplication {

	public static void main(String[] args) {

		SpringApplication.run(MCPServerApplication.class, args);
	}

	@Bean
	public ToolCallbackProvider dataTools(ToolsDefinitionService dataService) {

		return MethodToolCallbackProvider.builder().toolObjects(dataService).build();
	}

	/**
	 * 输出数据库查询确认服务端启动成功
	 */
	// @Resource
	// private MySQLDataService mySQLDataServiceImpl;
	//
	// @Bean
	// public Void run() {
	//
	// 	System.out.println(mySQLDataServiceImpl.getContextData("张三"));
	// 	return null;
	// }

}
