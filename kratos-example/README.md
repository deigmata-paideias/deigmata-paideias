# Kratos 拦截器

```shell
curl 127.0.0.1:8000/example/test
```

## http server 相关用法

### middleware 中间件 

https://go-kratos.dev/docs/component/transport/http/ 

`Middleware(m ...middleware.Middleware) ServerOption`

```go
var opts = []http.ServerOption{
    http.Middleware(
        recovery.Recovery(),
		// kratos 自带中间件注册
        logging.Server(logger),
    ),
}
```

### Filter 拦截器

自定义看代码吧。

期望看到的输出：
```text
2025/04/25 23:56:35 过滤器：ReqUrlFilter()#Request URL: /example/test
2025/04/25 23:56:35 中间件：RealIpFilterMiddleware()#Request Real IP: 127.0.0.1:2894
```

## Kratos jsonPb 空值不返回配置

kratos issue：https://github.com/go-kratos/kratos/issues/1952
ptotobuf 不添加 omitempty：https://github.com/golang/protobuf/issues/1371

```protobuf
message HelloReply {
  string message = 1;
  string noRespStringType = 2;
  float noRespFloatType = 3;
  repeated string noRespRepeatedStringType = 4;
}
```

正常返回，存在默认输出：

```json
{"message":"Hello test","noRespStringType":"","noRespFloatType":0,"noRespRepeatedStringType":[]}
```

### 1. ResponseEncoder(en EncodeResponseFunc) ServerOption

自定义 pb 的 json 序列化处理策略。

https://go-kratos.dev/docs/component/transport/http/#responseencoderen-encoderesponsefunc-serveroption

配置 kratos 的响应编码器: encoder 包：

看到如下输出：
```text
经过了 CustomProtoJson ..............................
```

### 2. 自定义响应编码器

不使用 kratos 的 json 序列化器，直接使用 json 序列化器
