package indi.yuluo.graph.customnode;

import com.alibaba.cloud.ai.graph.OverAllState;
import com.alibaba.cloud.ai.graph.action.NodeAction;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.HashMap;
import java.util.Map;

/**
 * @author yuluo
 * @author <a href="mailto:yuluo08290126@gmail.com">yuluo</a>
 */

public class RecordingNode implements NodeAction {

    private static final Logger logger = LoggerFactory.getLogger(RecordingNode.class);

    @Override
    public Map<String, Object> apply(OverAllState state) {

        String feedback = (String) state.value("classifier_output").get();

        Map<String, Object> updatedState = new HashMap<>();
        if (feedback.contains("positive")) {
            logger.info("Received positive feedback: {}", feedback);
            updatedState.put("solution", "Praise, no action taken.");
        }
        else {
            logger.info("Received negative feedback: {}", feedback);
            updatedState.put("solution", feedback);
        }

        return updatedState;
    }

}
