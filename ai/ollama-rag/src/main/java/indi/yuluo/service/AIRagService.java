package indi.yuluo.service;

import java.io.IOException;
import java.nio.charset.StandardCharsets;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import reactor.core.publisher.Flux;

import org.springframework.ai.chat.client.ChatClient;
import org.springframework.ai.chat.client.advisor.QuestionAnswerAdvisor;
import org.springframework.ai.chat.model.ChatModel;
import org.springframework.ai.vectorstore.SearchRequest;
import org.springframework.ai.vectorstore.VectorStore;
import org.springframework.ai.vectorstore.filter.FilterExpressionBuilder;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.core.io.Resource;
import org.springframework.stereotype.Service;

/**
 * @author yuluo
 * @author <a href="mailto:yuluo08290126@gmail.com">yuluo</a>
 */

@Service
public class AIRagService {

	private final Logger logger = LoggerFactory.getLogger(AIRagService.class);

	@Value("classpath:/prompts/system-qa.st")
	private Resource systemResource;

	private final ChatModel ragChatModel;

	private final ChatClient ragClient;

	private final VectorStore vectorStore;

	private static final String textField = "content";

	public AIRagService(
			ChatModel chatModel,
			// 使用 DashScope ChatModel
			// @Qualifier("dashscopeChatModel") ChatModel chatModel,
			VectorStore vectorStore
	) {

		this.ragChatModel = chatModel;
		this.vectorStore = vectorStore;

		this.ragClient = ChatClient.builder(ragChatModel)
				.defaultAdvisors((new QuestionAnswerAdvisor(vectorStore)))
				.build();
	}

	public Flux<String> retrieve(String prompt) {

		// Get the vector store prompt tmpl.
		String promptTemplate = getPromptTemplate(systemResource);

		// Enable hybrid search, both embedding and full text search
		SearchRequest searchRequest = SearchRequest.builder().
				topK(4)
				.similarityThresholdAll()
				.filterExpression(
						new FilterExpressionBuilder()
								.eq(textField, prompt).build()
				).build();

		// Build ChatClient with retrieval rerank advisor:
		// ChatClient runtimeChatClient = ChatClient.builder(chatModel)
		//		.defaultAdvisors(new RetrievalRerankAdvisor(vectorStore, rerankModel, searchRequest, promptTemplate, 0.1))
		//		.build();

		// Spring AI RetrievalAugmentationAdvisor
		// Advisor retrievalAugmentationAdvisor = RetrievalAugmentationAdvisor.builder()
		//		.queryTransformers(RewriteQueryTransformer.builder()
		//				.chatClientBuilder(ChatClient.builder(ragChatModel).build().mutate())
		//				.build())
		//		.documentRetriever(VectorStoreDocumentRetriever.builder()
		//				.similarityThreshold(0.50)
		//				.vectorStore(vectorStore)
		//				.build())
		//		.build();

		// Retrieve and llm generate
		// return ragClient.prompt()
		//		.advisors(retrievalAugmentationAdvisor)
		//		.user(prompt)
		//		.stream()
		//		.content();

		return ChatClient.builder(ragChatModel)
				.build().prompt()
				.advisors(new QuestionAnswerAdvisor(vectorStore, searchRequest, promptTemplate))
				.user(prompt)
				.stream()
				.content();
	}

	private String getPromptTemplate(Resource systemResource) {

		try {
			logger.info("Loading system resource: {}", systemResource.getURI());
			return systemResource.getContentAsString(StandardCharsets.UTF_8);
		}
		catch (IOException e) {
			throw new RuntimeException(e);
		}
	}

}
