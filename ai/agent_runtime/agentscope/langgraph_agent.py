from typing import TypedDict

from langgraph import graph, types
from agentscope_runtime.engine.agents.langgraph_agent import LangGraphAgent

# 定义状态
class State(TypedDict, total=False):
    id: str


# 定义节点函数
async def set_id(state: State):
    new_id = state.get("id")
    assert new_id is not None, "must set ID"
    return types.Command(update=State(id=new_id), goto="REVERSE_ID")


async def reverse_id(state: State):
    new_id = state.get("id")
    assert new_id is not None, "ID must be set before reversing"
    return types.Command(update=State(id=new_id[::-1]))


state_graph = graph.StateGraph(state_schema=State)
state_graph.add_node("SET_ID", set_id)
state_graph.add_node("REVERSE_ID", reverse_id)
state_graph.set_entry_point("SET_ID")
compiled_graph = state_graph.compile(name="ID Reversal")

# 声明 agent
langgraph_agent = LangGraphAgent(graph=compiled_graph)

print("✅ LangGraph agent created successfully")
