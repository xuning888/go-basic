package main

import (
	"github.com/gin-gonic/gin"
	"go-basic/webbook/internal/dal"
	"go-basic/webbook/internal/repository"
	"go-basic/webbook/internal/repository/cache"
	"go-basic/webbook/internal/repository/dao"
	"go-basic/webbook/internal/service"
	"go-basic/webbook/internal/service/sms/memory"
	"go-basic/webbook/internal/web"
	"log"
	"net/http"
)

func main() {
	server := initWebServer()

	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello world")
	})

	log.Fatal(server.Run(":8080"))
}

func initWebServer() *gin.Engine {
	// 初始化基础组件
	db := dal.InitDb()
	redis := dal.InitRedis()
	// 初始化repo
	userRepository := repository.NewUserRepository(dao.NewUserDAO(db), cache.NewRedisUserCache(redis))
	codeRepository := repository.NewCodeRepository(cache.NewRedisCodeCache(redis))
	// 初始化service
	smsService := memory.NewMemService()
	userService := service.NewUserService(userRepository)
	codeService := service.NewCodeService(codeRepository, smsService)
	userHandler := web.NewUserHandler(userService, codeService)
	return dal.InitWebServer(redis, userHandler)
}
