package server

import (
	"fmt"
	v1 "kratos-example/api/example/v1"
	"kratos-example/internal/conf"
	"kratos-example/internal/encoder"
	"kratos-example/internal/filter"
	"kratos-example/internal/middleware"
	"kratos-example/internal/service"
	sHttp "net/http"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.ExampleService, logger log.Logger) *http.Server {

	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
			middleware.RealIpFilterMiddleware(),
		),
		// Filter 用的不多似乎？
		http.Filter(
			filter.ReqUrlFilter,
		),
		// 注册自定义 json encoder
		// json 自带有判断 proto 和 原始 json 类型的判断，因此没有用到 json  里的
		// json:"name,omitempty"` 字段属性，
		// 强制 http 响应全部经过非 kratos 的默认编码器返回，就可以避免这个问题
		// 但是需要手动修改生成 proto go 文件中的 json 字段
		http.ResponseEncoder(encoder.CustomResponseEncoder),
	}

	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}

	srv := http.NewServer(opts...)

	// 实现了 go sdk 的 http HandlerFunc 接口
	// 可以接入 gin 的 handler
	srv.HandleFunc("/handler/func", func(writer sHttp.ResponseWriter, request *sHttp.Request) {
		// handlerFunc 用法
		fmt.Println("==================== handlerFunc ====================")
		fmt.Println(writer.Header())
		fmt.Println(request.URL)
	})

	// 测试输出
	endpoint, _ := srv.Endpoint()
	fmt.Printf("服务启动的 Endpoint：%s\n", endpoint.String())

	// 注册 http 服务
	v1.RegisterExampleHTTPServer(srv, greeter)

	return srv
}
