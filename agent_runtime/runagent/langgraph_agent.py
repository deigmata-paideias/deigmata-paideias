# langgraph_agent.py
from langgraph.graph import StateGraph, END
from typing import TypedDict, List

class ConversationState(TypedDict):
    messages: List[str]
    context: dict
    response: str


def understand_intent(state):
    # Analyze user intent
    return {"context": {"intent": "question"}}


def generate_response(state):
    # Generate appropriate response
    return {"response": f"You asked: {state['messages'][-1]}"}


# Build graph
workflow = StateGraph(ConversationState)
workflow.add_node("understand", understand_intent)
workflow.add_node("respond", generate_response)
workflow.add_edge("understand", "respond")
workflow.add_edge("respond", END)
workflow.set_entry_point("understand")

app = workflow.compile()

def invoke(input_data):
    result = app.invoke({"messages": [input_data["query"]]})
    return {"response": result["response"]}

if __name__ == "__main__":
    print(invoke({"query": "What is the capital of France?"}))
