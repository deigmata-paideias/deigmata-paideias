package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// curl http://localhost:8080/
// curl http://localhost:8080/ping

func main() {

	// 创建一个默认的路由引擎，无任何 middleware
	r := gin.New()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!")
	})

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// 将路由引擎作为 http.Handler 传递给 http.ListenAndServe()
	http.Handle("/", r)

	_ = r.Run(":8080")

}
