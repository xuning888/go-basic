package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	// 这个engine 在golang中有非常核心的地位，他负责路由转发和类似aop的操作
	server := gin.Default()

	// 静态路由
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello, gin")
	})

	// 参数路由
	server.GET("/users/:name", func(ctx *gin.Context) {
		name := ctx.Param("name")
		ctx.String(http.StatusOK, "hello, 这是参数路由, name:"+name)
	})

	// 这是通配符路由
	server.GET("/views/*.html", func(ctx *gin.Context) {
		path := ctx.Param(".html")
		ctx.String(http.StatusOK, "这是通配符路由, path:"+path)
	})

	// 查询参数
	server.GET("/order", func(ctx *gin.Context) {
		oid := ctx.Query("id")
		ctx.String(http.StatusOK, "这是查询参数"+oid)
	})

	log.Fatal(server.Run(":8080"))
}
