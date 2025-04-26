package middleware

import (
	"context"
	"log"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// RealIpFilterMiddleware is a middleware that logs the real IP address of the client
func RealIpFilterMiddleware() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				kind := tr.Kind().String()
				operation := tr.Operation()
				// 断言成 HTTP 的 Transport 可以拿到特殊信息
				if ht, ok := tr.(*http.Transport); ok {
					log.Printf("中间件：RealIpFilterMiddleware()#Request Real IP: %s", ht.Request().RemoteAddr)
				}
				log.Printf("其他信息：%s, %v", kind, operation)
			}

			// 继续调用下一个中间件
			return handler(ctx, req)
		}
	}
}
