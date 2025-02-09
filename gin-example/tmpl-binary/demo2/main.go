package main

import (
	"embed"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

// http://localhost:8080/public/assets/images/example.png
//
//go:embed assets/* templates/*
var f embed.FS

func main() {

	r := gin.Default()

	tmpl := template.Must(template.New("").ParseFS(
		f,
		"templates/*.tmpl",
		"templates/foo/*.tmpl"),
	)
	r.SetHTMLTemplate(tmpl)

	// example: /public/assets/images/example.png
	r.StaticFS("/public", http.FS(f))

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "bar.tmpl", gin.H{
			"title": "Gin Index Site",
		})
	})

	r.GET("/foo", func(c *gin.Context) {
		c.HTML(http.StatusOK, "bar.tmpl", gin.H{
			"title": "Gin Foo Site",
		})
	})

	r.GET("favicon.ico", func(c *gin.Context) {
		file, _ := f.ReadFile("assets/favicon.ico")
		c.Data(
			http.StatusOK,
			"img/x-icon",
			file,
		)
	})

	_ = r.Run(":8080")

}
