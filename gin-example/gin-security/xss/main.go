package main

import (
	"github.com/gin-gonic/gin"
	"html"
	"net/http"
)

// curl http://127.0.0.1:8080/input&input="你好"

func main() {

	r := gin.Default()

	r.GET("/input", func(c *gin.Context) {
		// 对用户输入内容进行转义，防止 xss 攻击
		input := c.Query("name")
		safeInput := html.EscapeString(input)
		c.JSON(http.StatusOK, gin.H{
			"input": safeInput,
		})
	})

	_ = r.Run(":8080")

}
