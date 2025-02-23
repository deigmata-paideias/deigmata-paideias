package indi.yuluo.controller;

import indi.yuluo.service.AIRagService;
import jakarta.annotation.Resource;
import jakarta.servlet.http.HttpServletResponse;
import reactor.core.publisher.Flux;

import org.springframework.util.StringUtils;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

/**
 * @author yuluo
 * @author <a href="mailto:yuluo08290126@gmail.com">yuluo</a>
 */

@RestController
@RequestMapping("/rag/ai")
public class AIRagController {

	@Resource
	public AIRagService aiRagService;

	@GetMapping("/chat/{prompt}")
	public Flux<String> chat(
			@PathVariable("prompt") String prompt,
			HttpServletResponse response
	) {

		response.setCharacterEncoding("UTF-8");

		if (!StringUtils.hasText(prompt)) {
			return Flux.just("prompt is null.");
		}

		return aiRagService.retrieve(prompt);
	}

}
