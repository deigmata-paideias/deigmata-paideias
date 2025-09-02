import asyncio
from agentscope_runtime.engine.deployers import LocalDeployManager

# 倒入 agentscope runner
from run import create_runner

async def deploy_agent(runner):
    # 创建部署管理器
    deploy_manager = LocalDeployManager(
        host="localhost",
        port=8090,
    )

    # 将智能体部署为流式服务
    deploy_result = await runner.deploy(
        deploy_manager=deploy_manager,
        endpoint_path="/process",
        stream=True,  # Enable streaming responses
    )
    print(f"🚀智能体部署在: {deploy_result}")
    print(f"🌐服务URL: {deploy_manager.service_url}")
    print(f"💚 健康检查: {deploy_manager.service_url}/health")

    return deploy_manager
    

async def run_agent():
    async with create_runner() as runner:
        deploy_mamager = await deploy_agent(runner)

        # 保持服务运行
        print("service is running...")
        try:
            # 保持运行直到中断
            while True:
                await asyncio.sleep(1)
        except KeyboardInterrupt:
            print("服务已经停止了...")

if __name__ == "__main__":
    asyncio.run(run_agent())
    