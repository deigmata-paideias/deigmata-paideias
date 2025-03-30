---
slug: spring-ai-alibaba-observations
title: Spring AI 可观测性
date: 2025-03-30 13:07:00
authors: yuluo
tags: [Spring AI observations, ZipKin]
keywords: [Spring AI observations, grafana, prometheus]
image: /img/ai/observations/1.png
---

<!-- truncate -->

Spring AI 通过集成 micrometer 来完成 Model 中观测能力。MicroMeter 本身支持对 otel，zipkin 等多种形式的数据格式的导出。此文中通过 zipkin 和 grafana 的方式演示 Spring AI 的可观测功能。

Github 代码地址：https://github.com/deigmata-paideias/deigmata-paideias/tree/main/ai/observations

## Spring AI 可观测 Example

### Zipkin 

> 也可以使用 jaeger 代替 zipkin。

#### 代码和环境准备

引入 micrometer 和 zipkin 依赖：

```xml
<dependencies>
    <dependency>
        <groupId>org.springframework.boot</groupId>
        <artifactId>spring-boot-starter-web</artifactId>
        <version>3.4.2</version>
    </dependency>

    <dependency>
        <groupId>com.alibaba.cloud.ai</groupId>
        <artifactId>spring-ai-alibaba-starter</artifactId>
        <version>1.0.0-M6.1</version>
    </dependency>

    <dependency>
        <groupId>org.springframework.boot</groupId>
        <artifactId>spring-boot-starter-actuator</artifactId>
        <version>3.4.2</version>
    </dependency>

    <dependency>
        <groupId>io.micrometer</groupId>
        <artifactId>micrometer-tracing-bridge-brave</artifactId>
        <version>1.5.0-M2</version>
        <exclusions>
            <exclusion>
                <artifactId>slf4j-api</artifactId>
                <groupId>org.slf4j</groupId>
            </exclusion>
        </exclusions>
    </dependency>

    <dependency>
        <groupId>io.zipkin.reporter2</groupId>
        <artifactId>zipkin-reporter-brave</artifactId>
        <version>3.4.3</version>
    </dependency>

</dependencies>
```

application.yml 配置

> 配置 Spring AI 可观测性的配置开关和 actuator 的 endpoint 相关配置。
>
> 并没有配置 apikey，spring ai alibaba starter 会默认在环境变量里查询相关 key。

```yml
 spring:
  ai:
	# 注意：观测数据可能存在敏感信息！
    chat:
      client:
        observations:
          # 记录调用输入
          include-input: true
      observations:
        # 记录 LLMs 输出
        include-completion: true
        # 记录 prompt
        include-prompt: true
        include-error-logging: true

# 部署时不配置，使用默认配置
management:
  endpoints:
    web:
      exposure:
        # 开放所有外部端点
        include: "*"
  endpoint:
    health:
      # 应用健康状态检查，携带详细新学年
      show-details: always
  tracing:
    sampling:
        # trace 采样信息，记录每个请求
        probability: 1.0
```

docker compose zipkin 

> 本例中将监控数据导入到 zipkin 

```yml
services:

  zipkin:
    image: 'openzipkin/zipkin:latest'
    container_name: "sa-zipkin"
    ports:
      - '9411:9411'
```

chat 应用

```java
@RestController
@RequestMapping("/chat")
public class ChatController {

	private final ChatClient client;

	private ChatController(ChatModel model) {

		this.client = ChatClient.builder(model)
				.build();
	}

	@GetMapping
	public Flux<String> chat(
			@RequestParam("prompt") String prompt,
			HttpServletResponse response
	) {

		response.setCharacterEncoding("UTF-8");
		return client.prompt()
				.user(prompt)
				.stream()
				.content();
	}
}
```

#### 效果演示

##### 访问接口

```shell
# requset
curl http://localhost:8080/chat?prompt=你好

# response
你好！今天过得怎么样？
```

##### zipkin 查看请求


![zipkin](1.png)

可以看到，请求的相关信息已经打印到了 zipkin 控制台。

### prometheus + Grafana

#### 代码和环境准备

pom.yml

> 将上面 zipkin 的依赖替换为：

```xml
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-web</artifactId>
    <version>3.4.2</version>
</dependency>

<dependency>
    <groupId>com.alibaba.cloud.ai</groupId>
    <artifactId>spring-ai-alibaba-starter</artifactId>
    <version>1.0.0-M6.1</version>
</dependency>

<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-actuator</artifactId>
    <version>3.4.2</version>
</dependency>

<dependency>
    <groupId>io.micrometer</groupId>
    <artifactId>micrometer-registry-prometheus</artifactId>
    <version>1.14.5</version>
</dependency>
```

application.yml 改为：

```yml
spring:
  application:
    name: "spring-ai-observations-example"

  ai:
    # 注意：观测数据可能存在敏感信息！
    chat:
      client:
        observations:
          # 记录调用输入
          include-input: true
      observations:
        # 记录 LLMs 输出
        include-completion: true
        # 记录 prompt
        include-prompt: true
        include-error-logging: true

management:
  endpoints:
    web:
      exposure:
        include: "*"
  endpoint:
    health:
      show-details: always
  tracing:
    sampling:
      probability: 1.0
  prometheus:
    metrics:
      export:
        enabled: true
  metrics:
    tags:
      application: ${spring.application.name}
```

部署 grafana 和 prometheus 和 docker-compose.yml

```yml
services:

  prometheus:
    container_name: "sa-prometheus"
    hostname: "sa-prometheus"
    image: prom/prometheus:latest
    volumes:
      - ./prometheus/prometheus.yml:/etc/prometheus/prometheus.yml
    ports:
      - "9090:9090"

  grafana:
    image: grafana/grafana
    hostname: "sa-grafana"
    container_name: "sa-grafana"
    ports:
      - "3000:3000"

networks:
  bridge:
    driver: bridge
```

prometheus 配置文件：<span style="color:#990000;">(在 docker 中启动时，需要注意 job targets 配置成宿主机的 ipv4 地址，否则会连接失败！)</span>

```yml
scrape_configs:
  - job_name: "prometheus"
    static_configs:
      - targets: ["localhost:9090"]
    # prometheus 上面是默认的数据

  # 抓取任务配置
  - job_name: 'policy-prometheus'
    # 抓取目标的路径,其中?替换为要监听的服务前缀端口
    # micrometer 端点默认没有前缀
    metrics_path: '/actuator/prometheus'
    static_configs:
      # 抓取任务的目标地址，可以使用数组的形式配置多个端口监听
      - targets: ["192.168.99.117:8080"]
```

#### 效果演示

1. 访问 http://localhost:8080/actuator/prometheus 即可看到 ai 应用暴露的指标信息；
2. 访问 prometheus  http://localhost:9090/targets?search= 可以看到 jobs 信息；
3. 访问 grafana ，在导入数据源和 dashboard 之后可以看到 metrics 面板。

