# Go-Zero example

go-zero 是一个集成了各种工程实践的 web 和 rpc 框架。
通过弹性设计保障了大并发服务端的稳定性，经受了充分的实战检验。

## Go-Zero 环境安装

### goctl

goctl 是 go-zero 的内置脚手架，是提升开发效率的一大利器，可以一键生成代码、文档、部署 k8s yaml、dockerfile 等。

```shell
go install github.com/zeromicro/go-zero/tools/goctl@latest
```

验证：

```shell
goctl --version
```

### etcd

go-zero 服务发现使用。

docker-compose 安装：

```yml
version: "3.3"

services:

  etcd:
    container_name: go-zero-etcd
    hostname: etcd
    image: bitnami/etcd:3
    deploy:
      replicas: 1
      restart_policy:
        condition: on-failure
    privileged: true
    environment:
      - "ETCD_ADVERTISE_CLIENT_URLS=http://0.0.0.0:2379"
      - "ETCD_LISTEN_CLIENT_URLS=http://0.0.0.0:2379"
      - "ETCD_LISTEN_PEER_URLS=http://0.0.0.0:2380"
      - "ETCD_INITIAL_ADVERTISE_PEER_URLS=http://0.0.0.0:2380"
      #参数指不用密码可以连接
      - "ALLOW_NONE_AUTHENTICATION=yes"
      - "ETCD_INITIAL_CLUSTER=node1=http://0.0.0.0:2380"
      - "ETCD_NAME=node1"
      - "ETCD_DATA_DIR=/opt/bitnami/etcd/data"
    ports:
      # 修改了端口，2379 在 windows 是保留端口
      - 8079:2379
      - 8080:2380
    networks:
      - go-zero-example

  # etcd ui console.
  etcdkeeper:
    image: deltaprojects/etcdkeeper
    container_name: go-zero-etcdkeeper
    ports:
      - 8088:8080
    networks:
      - go-zero-example

networks:
  go-zero-example:
```

启动之后，需要在 etcekeeper 设置 etcdserver 地址：etcd:2379

### Protoc

protoc 是一个用于生成代码的工具，它可以根据 proto 文件生成C++、Java、Python、Go、PHP 等多重语言的代码，
而 gRPC 的代码生成还依赖 protoc-gen-go，protoc-gen-go-grpc 插件来配合生成 Go 语言的 gRPC 代码。

使用 goctl 安装：

```shell
goctl env check --install --verbose --force
```

### grpcurl

类似 curl，对 proto 接口发起调用。

```shell
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

### Go-Zero 安装

需要先创建一个项目，在项目中执行。

```shell
go get -u github.com/zeromicro/go-zero@latest
```

go 需要设置代理

```shell
$ go env -w GOPROXY=https://goproxy.cn,direct
$ go env GOPROXY
https://goproxy.cn,direct
```

### goctl-intellij 

goctl-intellij 是 go-zero api 描述语言的 intellij 编辑器插件，
支持 api 描述语言高亮、语法检测、快速提示、创建模板特性。

Intellij plugin 搜索安装。


## Go-Zero 快速开始

需求场景：

> 视频微服务提供一个 http 接口，用户查询一个视频信息，并且将关联用户 id 的用户名也查出来。

### 用户微服务模块

新建 `user/rpc/user.proto` 文件，写入以下内容：

```protobuf
syntax = "proto3";

package user;

option go_package = "./user";

message IdRequest {
  string id = 1;
}

message UserResponse {
  string id = 1;
  string name = 2;
  string gender = 3;
}

service User {
  rpc getUser(IdRequest) returns (UserResponse);
}

// 根目录执行：
// goctl rpc protoc user/rpc/user.proto --go_out=user/rpc/types --go-grpc_out=user/rpc/types --zrpc_out=user/rpc
```

最后在根目录执行：`goctl rpc protoc user/rpc/user.proto --go_out=user/rpc/types --go-grpc_out=user/rpc/types --zrpc_out=user/rpc
会生成一系列 user 的项目代码，

在 /user/rpc/internal/logic/getUserLogic.go 中实现业务逻辑：

```go
func (l *GetUserLogic) GetUser(in *user.IdRequest) (*user.UserResponse, error) {

	return &user.UserResponse{
		Id: "1234",
		Name: "yuluo",
		Gender: true,
}, nil
}
```

之后修改 `user/rpc/etc/user.yaml` 中 etcd 的地址为 docker 部署的 etcd 的地址：`127.0.0.1:8079`。

```yaml
Name: user.rpc
# user 微服务端口
ListenOn: 0.0.0.0:18080
Etcd:
  Hosts:
  # etcd 地址
  - 127.0.0.1:8079
  # etcd 服务发现的名字
  Key: user.rpc
```

启动服务验证：

```shell
# 下载项目依赖
go mod tidy

# 启动服务
go run user.go
```

### 视频微服务模块

根目录创建 video/api/video.api 文件：

```api
type (
	VideoReq {
		Id string `path: "id"`
	}
	VideoRes {
		Id   string `json:"id"`
		Name string `json:"name"`
	}
)

service VideoService {
	@handler getVideo
	get /api/video/:id (VideoReq) returns (VideoRes)
}

// goctl api go -api video/api/video.api -dir video/api
```

生成 video 代码：

```shell
goctl api go -api video/api/video.api -dir video/api
```

在 video/api/internal/logic/getvideologic.go 中实现业务逻辑，调用 rpc 服务：

在 video/api/internal/config/config.go 中配置 userRpc：

```go
type Config struct {
	rest.RestConf
	userRpc zrpc.RpcClientConf
}
```

在 video/api/internal/svc/servicecontext.go 完善 userRpc：

```go
type ServiceContext struct {
	Config  config.Config
	UserRpc userclient.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		UserRpc: userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}
```

在 video/api/internal/handler/getvideohandler.go 中调用 userRpc：

```go
func (l *GetVideoLogic) GetVideo(req *types.VideoReq) (resp *types.VideoRes, err error) {

	userl, err := l.svcCtx.UserRpc.GetUser(l.ctx, &user.IdRequest{
		Id: "1234",
	})

	if err != nil {
		return nil, err
	}

	return &types.VideoRes{
		Id:   req.Id,
		Name: userl.Name,
	}, nil
}
```

在 video/api/etc/video.yaml 中配置 user 微服务 地址：

```yaml
Name: VideoService
Host: 0.0.0.0
Port: 8888

UserRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:8079
    key:  user.rpc
```

启动：

```shell
cd video/api
go mod tidy
go run videoservice.go
```

访问测试：

```shell
$ curl 127.0.0.1:8888/api/video/1234

{"id":"1234","name":"yuluo"}
```
