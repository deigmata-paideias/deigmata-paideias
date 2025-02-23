package indi.yuluo.core;

import java.io.IOException;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import co.elastic.clients.elasticsearch.ElasticsearchClient;
import co.elastic.clients.elasticsearch._types.mapping.DenseVectorProperty;
import co.elastic.clients.elasticsearch._types.mapping.KeywordProperty;
import co.elastic.clients.elasticsearch._types.mapping.ObjectProperty;
import co.elastic.clients.elasticsearch._types.mapping.Property;
import co.elastic.clients.elasticsearch._types.mapping.TextProperty;
import co.elastic.clients.elasticsearch._types.mapping.TypeMapping;
import co.elastic.clients.elasticsearch.indices.CreateIndexResponse;
import co.elastic.clients.elasticsearch.indices.IndexSettings;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import org.springframework.ai.autoconfigure.vectorstore.elasticsearch.ElasticsearchVectorStoreProperties;
import org.springframework.ai.document.Document;
import org.springframework.ai.document.DocumentReader;
import org.springframework.ai.reader.pdf.PagePdfDocumentReader;
import org.springframework.ai.transformer.splitter.TokenTextSplitter;
import org.springframework.ai.vectorstore.VectorStore;
import org.springframework.boot.ApplicationArguments;
import org.springframework.boot.ApplicationRunner;
import org.springframework.core.io.Resource;
import org.springframework.core.io.ResourceLoader;
import org.springframework.core.io.support.ResourcePatternResolver;
import org.springframework.stereotype.Component;
import org.springframework.util.StringUtils;

/**
 * @author yuluo
 * @author <a href="mailto:yuluo08290126@gmail.com">yuluo</a>
 */

@Component
public class ApplicationInit implements ApplicationRunner {

	private final Logger logger = LoggerFactory.getLogger(ApplicationInit.class);

	private final ResourceLoader resourceLoader;

	private final VectorStore vectorStore;

	private final ElasticsearchClient elasticsearchClient;

	private final ElasticsearchVectorStoreProperties options;

	private static final String textField = "content";

	private static final String vectorField = "embedding";

	public ApplicationInit(
			ResourceLoader resourceLoader,
			VectorStore vectorStore,
			ElasticsearchClient elasticsearchClient,
			ElasticsearchVectorStoreProperties options
	) {

		this.resourceLoader = resourceLoader;
		this.vectorStore = vectorStore;
		this.elasticsearchClient = elasticsearchClient;
		this.options = options;
	}

	@Override
	public void run(ApplicationArguments args) {

		// 1. load pdf resources.
		List<Resource> pdfResources = loadPdfResources();

		// 2. parse pdf resources to Documents.
		List<Document> documents = parsePdfResource(pdfResources);

		// 3. import to ES.
		importToES(documents);

		logger.info("RAG application init finished");
	}

	private List<Resource> loadPdfResources() {

		List<Resource> pdfResources = new ArrayList<>();

		try {

			logger.info("加载 PDF 资源=================================");

			ResourcePatternResolver resolver = (ResourcePatternResolver) resourceLoader;
			Resource[] resources = resolver.getResources("classpath:data" + "/*.pdf");

			for (Resource resource : resources) {
				if (resource.exists()) {
					pdfResources.add(resource);
				}
			}

			logger.info("加载 PDF 资源完成=================================");
		}
		catch (Exception e) {
			throw new RuntimeException("Failed to load PDF resources", e);
		}

		return pdfResources;
	}

	private List<Document> parsePdfResource(List<Resource> pdfResources) {

		List<Document> resList = new ArrayList<>();

		logger.info("开始解析 PDF 资源=================================");

		for (Resource springAiResource : pdfResources) {

			// 1. parse document
			DocumentReader reader = new PagePdfDocumentReader(springAiResource);
			List<Document> documents = reader.get();
			logger.info("{} documents loaded", documents.size());

			// 2. split trunks
			List<Document> splitDocuments = new TokenTextSplitter().apply(documents);
			logger.info("{} documents split", splitDocuments.size());

			// 3. add res list
			resList.addAll(splitDocuments);
		}

		logger.info("解析 PDF 资源完成=================================");

		return resList;
	}

	private void importToES(List<Document> documents) {

		logger.info("开始导入数据到 ES =================================");

		logger.info("create embedding and save to vector store");
		createIndexIfNotExists();
		vectorStore.add(documents);

		logger.info("导入数据到 ES 完成=================================");
	}

	private void createIndexIfNotExists() {

		try {
			String indexName = options.getIndexName();
			Integer dimsLength = options.getDimensions();

			if (!StringUtils.hasText(indexName)) {
				throw new IllegalArgumentException("Elastic search index name must be provided");
			}

			boolean exists = elasticsearchClient.indices().exists(idx -> idx.index(indexName)).value();
			if (exists) {
				logger.debug("Index {} already exists. Skipping creation.", indexName);
				return;
			}

			String similarityAlgo = options.getSimilarity().name();
			IndexSettings indexSettings = IndexSettings
					.of(settings -> settings.numberOfShards(String.valueOf(1)).numberOfReplicas(String.valueOf(1)));

			Map<String, Property> properties = new HashMap<>();
			properties.put(vectorField, Property.of(property -> property.denseVector(
					DenseVectorProperty.of(dense -> dense.index(true).dims(dimsLength).similarity(similarityAlgo)))));
			properties.put(textField, Property.of(property -> property.text(TextProperty.of(t -> t))));

			Map<String, Property> metadata = new HashMap<>();
			metadata.put("ref_doc_id", Property.of(property -> property.keyword(KeywordProperty.of(k -> k))));

			properties.put("metadata",
					Property.of(property -> property.object(ObjectProperty.of(op -> op.properties(metadata)))));

			CreateIndexResponse indexResponse = elasticsearchClient.indices()
					.create(createIndexBuilder -> createIndexBuilder.index(indexName)
							.settings(indexSettings)
							.mappings(TypeMapping.of(mappings -> mappings.properties(properties))));

			if (!indexResponse.acknowledged()) {
				throw new RuntimeException("failed to create index");
			}

			logger.info("create elasticsearch index {} successfully", indexName);
		}
		catch (IOException e) {
			logger.error("failed to create index", e);
			throw new RuntimeException(e);
		}
	}

}
