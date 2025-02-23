package indi.yuluo.configuration;

import org.springframework.boot.autoconfigure.AutoConfiguration;

/**
 * @author yuluo
 * @author <a href="mailto:yuluo08290126@gmail.com">yuluo</a>
 *
 * 也可在配置文件中配置
 */

@AutoConfiguration
public class ESConfiguration {

//	@Bean
//	public RestClient restClient() {
//
//		return RestClient.builder(
//				new HttpHost(
//						"127.0.0.1",
//						9200,
//						"http")
//				).build();
//	}
//
//	/**
//	 * 自定义 ES VectorStore 解决:
//	 * [1:15386] failed to parse: The [dense_vector] field [embedding] in doc
//	 * [document with id '464d6d6d-e685-4432-b0c2-699af3ff3bea'] has
//	 * a different number of dimensions [768] than defined in the mapping [1536]
//	 */
//	@Bean
//	public VectorStore vectorStore(RestClient restClient, EmbeddingModel embeddingModel) {
//
//		ElasticsearchVectorStoreOptions options = new ElasticsearchVectorStoreOptions();
//		options.setIndexName("ai-rag-index");
//		options.setSimilarity(SimilarityFunction.cosine);
//
//		// 和向量模型输出的向量维度一致
//		options.setDimensions(768);
//
//		return ElasticsearchVectorStore.builder(restClient, embeddingModel)
//				.options(options)
//				.initializeSchema(true)
//				.batchingStrategy(new TokenCountBatchingStrategy())
//				.build();
//	}

}
