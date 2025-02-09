# Gin 和 gRPC 集成使用

Gin 作为 HTTP 请求接口，之后通过 gRPC 客户端与 gRPC 服务建立通信。

## 安装 gRPC plugins

```shell
go install google.golang.org/protobuf/cmd/protoc-gen-go
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

更新环境变量，以便于找到之前安装的工具

```shell
export PATH="$PATH:$(go env GOPATH)/bin"
```

## 生成 go 的 gRPC 代码

```shell
protoc --go_out=gen --go_opt=paths=source_relative \
  --go-grpc_out=gen --go-grpc_opt=paths=source_relative \
  -I=$PWD proto/v1/helloworld.proto
```

## 运行

```shell
go run gRPC/server.go

go run gin/main.go
```

## 测试

```shell
curl 'http://localhost:8080/rest/n/gin'
```

使用 grpcurl 直接测试 gRPC 服务

```shell
grpcurl -d '{"name": "gin"}' \
  -plaintext localhost:50051 helloworld.v1.Greeter/SayHello
```
