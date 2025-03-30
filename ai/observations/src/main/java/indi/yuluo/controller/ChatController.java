package indi.yuluo.controller;

import jakarta.servlet.http.HttpServletResponse;
import reactor.core.publisher.Flux;

import org.springframework.ai.chat.client.ChatClient;
import org.springframework.ai.chat.model.ChatModel;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

/**
 * @author yuluo
 * @author <a href="mailto:yuluo08290126@gmail.com">yuluo</a>
 */

@RestController
@RequestMapping("/chat")
public class ChatController {

	private final ChatClient client;

	private ChatController(ChatModel model) {

		this.client = ChatClient.builder(model)
				.build();
	}

	@GetMapping
	public Flux<String> chat(
			@RequestParam("prompt") String prompt,
			HttpServletResponse response
	) {

		response.setCharacterEncoding("UTF-8");
		return client.prompt()
				.user(prompt)
				.stream()
				.content();
	}

}
