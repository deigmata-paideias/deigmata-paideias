package indi.yuluo;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

/**
 * @author yuluo
 * @author <a href="mailto:yuluo08290126@gmail.com">yuluo</a>
 */

@SpringBootApplication
public class MCPClientApplication {

	public static void main(String[] args) {

		SpringApplication.run(MCPClientApplication.class, args);
	}

	// @Bean
	// public CommandLineRunner predefinedQuestions(ChatClient.Builder chatClientBuilder, ToolCallbackProvider tools,
	// 		ConfigurableApplicationContext context) {
	//
	// 	return args -> {
	//
	// 		var chatClient = chatClientBuilder
	// 				.defaultTools(tools)
	// 				.build();
	//
	// 		System.out.println("\n>>> QUESTION: " + "帮我查询一下张三的信息");
	// 		System.out.println("\n>>> ASSISTANT: " + chatClient.prompt("帮我查询一下张三的信息").call().content());
	//
	// 		context.close();
	// 	};
	// }

}
