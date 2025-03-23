package indi.yuluo;

import java.util.Arrays;

import jakarta.servlet.http.HttpServletResponse;
import reactor.core.publisher.Flux;

import org.springframework.ai.chat.client.ChatClient;
import org.springframework.ai.chat.model.ChatModel;
import org.springframework.ai.model.function.FunctionCallback;
import org.springframework.ai.tool.ToolCallbackProvider;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

/**
 * @author yuluo
 * @author <a href="mailto:yuluo08290126@gmail.com">yuluo</a>
 */

@RestController
@RequestMapping("/mcp")
public class MCPClientController {

	private final ChatClient client;

	public MCPClientController(
			ChatModel chatModel,
			ToolCallbackProvider tools
	) {
		Arrays.stream(tools.getToolCallbacks()).map(FunctionCallback::getName).forEach(System.out::println);

		this.client = ChatClient.builder(chatModel)
				.defaultTools(tools)
				.build();
	}

	@GetMapping("/chat")
	public Flux<String> chat(
			@RequestParam("prompt") String prompt,
			HttpServletResponse response
	) {

		response.setCharacterEncoding("UTF-8");

		return client.prompt().user(prompt).stream().content();
	}

}
