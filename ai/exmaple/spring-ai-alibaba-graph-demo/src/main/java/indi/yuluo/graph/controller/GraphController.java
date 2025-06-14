package indi.yuluo.graph.controller;

import com.alibaba.cloud.ai.graph.CompiledGraph;
import com.alibaba.cloud.ai.graph.StateGraph;
import com.alibaba.cloud.ai.graph.exception.GraphStateException;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import java.util.Map;

/**
 * @author yuluo
 * @author <a href="mailto:yuluo08290126@gmail.com">yuluo</a>
 */

@RestController
@RequestMapping("/graph/demo")
public class GraphController {

    private final CompiledGraph compiledGraph;

    public GraphController(@Qualifier("workflowGraph") StateGraph stateGraph) throws GraphStateException {

        this.compiledGraph = stateGraph.compile();
    }

    @GetMapping("/chat")
    public String simpleChat(@RequestParam("query") String query) {

        return compiledGraph.invoke(Map.of("input", query)).flatMap(input -> input.value("solution")).get().toString();
    }

}
