package indi.yuluo.graph.config;

import com.alibaba.cloud.ai.graph.GraphRepresentation;
import com.alibaba.cloud.ai.graph.OverAllState;
import com.alibaba.cloud.ai.graph.OverAllStateFactory;
import com.alibaba.cloud.ai.graph.StateGraph;
import com.alibaba.cloud.ai.graph.action.EdgeAction;
import com.alibaba.cloud.ai.graph.exception.GraphStateException;
import com.alibaba.cloud.ai.graph.node.QuestionClassifierNode;
import com.alibaba.cloud.ai.graph.state.strategy.ReplaceStrategy;
import indi.yuluo.graph.customnode.RecordingNode;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.ai.chat.client.ChatClient;
import org.springframework.ai.chat.client.advisor.SimpleLoggerAdvisor;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import java.util.HashMap;
import java.util.List;
import java.util.Map;

import static com.alibaba.cloud.ai.graph.StateGraph.END;
import static com.alibaba.cloud.ai.graph.StateGraph.START;
import static com.alibaba.cloud.ai.graph.action.AsyncEdgeAction.edge_async;
import static com.alibaba.cloud.ai.graph.action.AsyncNodeAction.node_async;

/**
 * Graph Demo：首先判断评价正负，其次细分负面问题，最后输出处理方案。
 *
 * @author yuluo
 * @author <a href="mailto:yuluo08290126@gmail.com">yuluo</a>
 */

@Configuration
public class GraphAutoConfiguration {

    private static final Logger logger = LoggerFactory.getLogger(GraphAutoConfiguration.class);

    /**
     * 定义一个工作流 StateGraph Bean.
     */
    @Bean
    public StateGraph workflowGraph(ChatClient.Builder builder) throws GraphStateException {

        // LLMs Bean
        ChatClient chatClient = builder.defaultAdvisors(new SimpleLoggerAdvisor()).build();

        // 定义一个 OverAllStateFactory，用于在每次执行工作流时创建初始的全局状态对象。通过注册若干 Key 及其更新策略来管理上下文数据
        // 注册三个状态 key 分别为
        // 1. input：用户输入的文本
        // 2. classifier_output：分类器的输出结果
        // 3. solution：最终输出结论
        // 使用 ReplaceStrategy（每次写入替换旧值）策略处理上下文状态对象中的数据，用于在节点中传递数据
        OverAllStateFactory stateFactory = () -> {
            OverAllState state = new OverAllState();
            state.registerKeyAndStrategy("input", new ReplaceStrategy());
            state.registerKeyAndStrategy("classifier_output", new ReplaceStrategy());
            state.registerKeyAndStrategy("solution", new ReplaceStrategy());
            return state;
        };

        // 创建 workflows 节点
        // 使用 Graph 框架预定义的 QuestionClassifierNode 来处理文本分类任务

        // 评价正负分类节点
        QuestionClassifierNode feedbackClassifier = QuestionClassifierNode.builder()
                .chatClient(chatClient)
                .inputTextKey("input")
                .categories(List.of("positive feedback", "negative feedback"))
                .classificationInstructions(
                        List.of("Try to understand the user's feeling when he/she is giving the feedback."))
                .build();

        // 负面评价具体问题分类节点
        QuestionClassifierNode specificQuestionClassifier = QuestionClassifierNode.builder()
                .chatClient(chatClient)
                .inputTextKey("input")
                .categories(List.of("after-sale service", "transportation", "product quality", "others"))
                .classificationInstructions(List
                        .of("What kind of service or help the customer is trying to get from us? Classify the question based on your understanding."))
                .build();

        // 编排 Node 节点，使用 StateGraph 的 API，将上述节点加入图中，并设置节点间的跳转关系
        // 首先将节点注册到图，并使用 node_async(...) 将每个 NodeAction 包装为异步节点执行（提高吞吐或防止阻塞，具体实现框架已封装）
        StateGraph stateGraph = new StateGraph("Consumer Service Workflow Demo", stateFactory)

                // 定义节点
                .addNode("feedback_classifier", node_async(feedbackClassifier))
                .addNode("specific_question_classifier", node_async(specificQuestionClassifier))
                .addNode("recorder", node_async(new RecordingNode()))

                // 定义边（流程顺序）
                .addEdge(START, "feedback_classifier")
                .addConditionalEdges("feedback_classifier",
                        edge_async(new FeedbackQuestionDispatcher()),
                        Map.of("positive", "recorder", "negative", "specific_question_classifier"))
                .addConditionalEdges("specific_question_classifier",
                        edge_async(new SpecificQuestionDispatcher()),
                        Map.of("after-sale", "recorder", "transportation", "recorder", "quality", "recorder", "others",
                                "recorder"))

                // 图的结束节点
                .addEdge("recorder", END);

        GraphRepresentation graphRepresentation = stateGraph.getGraph(GraphRepresentation.Type.PLANTUML,
                "workflow graph");

        System.out.println("\n\n");
        System.out.println(graphRepresentation.content());
        System.out.println("\n\n");

        return stateGraph;
    }

    public static class FeedbackQuestionDispatcher implements EdgeAction {

        @Override
        public String apply(OverAllState state) {

            String classifierOutput = (String) state.value("classifier_output").orElse("");
            logger.info("classifierOutput: {}", classifierOutput);

            if (classifierOutput.contains("positive")) {
                return "positive";
            }
            return "negative";
        }

    }

    public static class SpecificQuestionDispatcher implements EdgeAction {

        @Override
        public String apply(OverAllState state) {

            String classifierOutput = (String) state.value("classifier_output").orElse("");
            logger.info("classifierOutput: {}", classifierOutput);

            Map<String, String> classifierMap = new HashMap<>();
            classifierMap.put("after-sale", "after-sale");
            classifierMap.put("quality", "quality");
            classifierMap.put("transportation", "transportation");

            for (Map.Entry<String, String> entry : classifierMap.entrySet()) {
                if (classifierOutput.contains(entry.getKey())) {
                    return entry.getValue();
                }
            }

            return "others";
        }

    }

}
