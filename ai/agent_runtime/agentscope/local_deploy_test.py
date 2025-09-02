import json
import requests

def test_deployed_agent():
    # 准备测试负载
    payload = {
        "input": [
            {
                "role": "user",
                "content": [
                    {"type": "text", "text": json.dumps({"id": "test-runtime-id"})},
                ],
            },
        ],
        "session_id": "test_session_001",
        "user_id": "test_user_001",
    }

    print("🧪 测试部署的智能体...")

    # 测试流式响应
    try:
        response = requests.post(
            "http://localhost:8090/process",
            json=payload,
            stream=True,
            timeout=30,
        )

        print("📡 流式响应:")
        for line in response.iter_lines():
            if line:
                print(f"{line.decode('utf-8')}")
        print("✅ 流式测试完成")
    except requests.exceptions.RequestException as e:
        print(f"❌ 流式测试失败: {e}")

    # 测试JSON响应（如果可用）
    try:
        response = requests.post(
            "http://localhost:8090/process",
            json=payload,
            timeout=30,
        )

        if response.status_code == 200:
            print(f"📄 JSON响应: {response.content}")
            print("✅ JSON测试完成")
        else:
            print(f"⚠️ JSON端点返回状态: {response.status_code}")

    except requests.exceptions.RequestException as e:
        print(f"ℹ️ JSON端点不可用或失败: {e}")


# Run the test
test_deployed_agent()
