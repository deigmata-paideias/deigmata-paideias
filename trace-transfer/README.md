#  异构网络 trace 追踪 Demo

调用链：

spring cloud gateway -> user-service -> |openfeign| -> kratos.

traceId 在 gateway 中被设置，从 gw -> user-svc -> kratos 一直透传。

## 服务启动并调用

先启动下游服务 kratos：

启动上游调用者 Java Spring Boot 应用 user-service；

启动网关，接受请求。

## zipkin 观测

gateway 和 user-serivce 都配置 zipkin，使用 docker 打开 zipkin 即可观测追踪指标。
