# Kratos Example

#### 项目创建

```shell
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest

kratos new kratos-example

cd kratos-example

# 初始化换进
make init
```

#### 项目结构

```markdown
kratos-example
├── Dockerfile		# license
├── LICENSE		
├── Makefile		# makefilw 指令
├── go.mod
├── go.sum
├── internal		# 业务逻辑代码
│   ├── biz			# 业务逻辑
│   ├── conf		# 跟目录 config 配置的解析
│   ├── data		# biz 的 repo 接口的实现，数据层
│   ├── server		# grpc 和 http 的 server
│   └── service		# 实现了 api 定义的服务层，类似 DDD 的 application 层，处理 DTO 到 biz 领域实体的转换(DTO -> DO)
├── openapi.yaml	 # openapi 文档
├── api				# v1 是 api 版本号 proto 文件是接口定义，http.go 和 grpc.go 是生成的，不需要修改。http web 服务，grpc 微服务 rpc 调用
├── configs			# 项目配置，timeout 地址等
├── cmd				# 可执行的 main 入口，其中 wire 是依赖注入组件
└── third_party		 # proto 依赖的东西，make api 编译的时候会一起编译
```

#### 接口开发

现在生成的 kratos 项目中修改 api 下面的二级目录为项目名

之后将 greeter.proto 改为项目名.proto. error_reason.proto 也修改下

运行 make api，windows 运行

```shell
 protoc --proto_path=api --proto_path=third_party --go_out=paths=source_relative:api --go-http_out=paths=source_relative:api --go-grpc_out=paths=source_relative:api --openapi_out=fq_schema_naming=true,default_response=false:. api/example/v1/*.proto
```

修改 service 和 server、biz、data 的相关 proto 代码，按照顺序依次修改：

service 实现 api，serivce 负责和 biz 进行业务交互，biz 和 data 进行数据交互，server 启动 http 和 grpc server。

更改完成代码之后：

```shell
cd /cmd/kratos-example/

# 重新维护依赖注入关系
wire
```

启动

```shell
kratos run
```

访问修改之后的 http://localhost:8000/example/yuluo 接口，成功

#### 接入数据库

数据库相关的代码都放在 data 下面。

使用 wire 在外部初始化完成 db client 之后，注入到  data 中

```go
func NewDB(c *conf.Data) *gorm.DB {

	dsn := "root:082916@tcp(127.0.0.1:3306)/data?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// 不建议在生产环境上使用
	if err = db.AutoMigrate(); err != nil {
		panic(err)
	}

	return db
}
```

在 data 中注入 db 

```go
type Data struct {
	db *gorm.DB
}

// NewData .
func NewData(c *conf.Data, logger log.Logger, db *gorm.DB) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{db: db}, cleanup, nil
}
```

在  wire 中加入

```go
var ProviderSet = wire.NewSet(NewData, NewDB, NewExampleRepo)
```

随后更新 wire ，观察到 NewDB 已经被 wire 管理注入。

使用配置文件管理 mysql dsn:

```yml
data:
  database:
    dsn: root:082916@tcp(127.0.0.1:3306)/data?charset=utf8mb4&parseTime=True&loc=Local
```

更新 config/config.pb.proto

```proto
message Data {
  message Database {
    string driver = 1;
  }
//  暂时不用 redis 注释调
//  message Redis {
//    string network = 1;
//    string addr = 2;
//    google.protobuf.Duration read_timeout = 3;
//    google.protobuf.Duration write_timeout = 4;
//  }
  Database database = 1;
//  Redis redis = 2;
}
```

执行 make config 更新 config.pb.go

win 执行 ` protoc --proto_path=internal/conf --proto_path=third_party --go_out=paths=source_relative:internal/conf ./internal/conf/*.proto`

使用 kratos run 验证