package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	rc "github.com/redis/go-redis/v9"
	"go-basic/webbook/config"
	"go-basic/webbook/internal/repository"
	"go-basic/webbook/internal/repository/dao"
	"go-basic/webbook/internal/service"
	"go-basic/webbook/internal/web"
	"go-basic/webbook/internal/web/middleware"
	"go-basic/webbook/pkg/middlewares/ratelimit"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	db := initDB()
	server := initWebServer()
	userHandler := initUser(db)
	userHandler.RegisterRoutes(server)

	/*	server := gin.Default()*/

	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "你好， 你来了")
	})

	log.Fatal(server.Run(":8081"))
}

// initWebServer 初始化 webServer
func initWebServer() *gin.Engine {
	server := gin.Default()

	// 跨域配置
	server.Use(cors.New(cors.Config{
		// AllowMethods: []string{"GET", "POST", "PUT"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"x-jwt-token"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "youcompany.com")
		},
		MaxAge: 12 * time.Hour,
	}))

	// 限流配置
	if config.Config.EnableLimit {
		redisClient := rc.NewClient(&rc.Options{
			Addr: config.Config.Redis.Add,
		})
		server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())
	}

	// 基于内存的存储 session 的实现
	// store := memstore.NewStore([]byte("sBvCQd28JynD7NWi"), []byte("DUVmChM4T3cAlNR8"))

	// 设置 选择session的存储方式, 默认用cookie
	// store := cookie.NewStore([]byte("secret"))

	// 基于redis 来存储session
	store, err := redis.NewStore(16, "tcp", config.Config.Redis.Add, "",
		[]byte("MxQP9pSI6BzUL9XVSZrdSeJm6Jbhw42z"), []byte("0XCRv2ip2KMbnId8hT8UowhPc6yiTrhK"))
	if err != nil {
		panic(err)
	}
	// 设置session的存储方式
	server.Use(sessions.Sessions("mysessions", store))

	// 校验 session 的 middleware
	server.Use(middleware.NewLoginJWTMiddlewareBuilder().
		Ignore("/users/signup").
		Ignore("/users/login").
		Build())

	return server
}

func initUser(db *gorm.DB) *web.UserHandler {
	userDAO := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(userDAO)
	userService := service.NewUserService(repo)
	return web.NewUserHandler(userService)
}

// initDB 初始化数据库
func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.Db.DSN))
	if err != nil {
		// 一旦初始化过程出错, 应用就不要启动
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
