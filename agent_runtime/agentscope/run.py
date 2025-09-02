import asyncio
import json
from agentscope_runtime.engine.schemas.agent_schemas import (
    MessageType,
    RunStatus,
    AgentRequest,
)
from agentscope_runtime.engine.services.context_manager import (
    ContextManager
)
from agentscope_runtime.engine import Runner

from contextlib import asynccontextmanager

# 引入 agent 
from langgraph_agent import langgraph_agent

@asynccontextmanager
async def create_runner():
    async with ContextManager() as context_manager:
        runner = Runner(
            agent=langgraph_agent,
            context_manager=context_manager,
        )
        print("✅ Runner创建成功")
        yield runner

async def interact_with_agent(runner):
    # Create a request
    request = AgentRequest(
        input=[
            {
                "role": "user",
                "content": [
                    {
                        "type": "text",
                        "text": json.dumps({"id": "this is run id"})
                    },
                ],
            },
        ],
    )

    # 流式获取响应
    print("🤖 智能体正在处理您的请求...")
    all_result = ""
    async for message in runner.stream_query(request=request):
        # Check if this is a completed message
        if (
            message.object == "message"
            and MessageType.MESSAGE == message.type
            and RunStatus.Completed == message.status
        ):
            all_result = message.content[0].text

    print(f"📝智能体回复: {all_result}")
    return all_result

async def main():
    async with create_runner() as runner:
        await interact_with_agent(runner)

if __name__ == "__main__":
    asyncio.run(main())
