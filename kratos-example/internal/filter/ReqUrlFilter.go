package filter

import (
	"log"
	"net/http"
)

func ReqUrlFilter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 记录请求的 URL
		log.Printf("过滤器：ReqUrlFilter()#Request URL: %s", r.URL.String())

		// 调用下一个处理器
		next.ServeHTTP(w, r)
	})
}
