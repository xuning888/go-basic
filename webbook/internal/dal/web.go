package dal

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	rc "github.com/redis/go-redis/v9"
	"go-basic/webbook/config"
	"go-basic/webbook/internal/web"
	"go-basic/webbook/internal/web/middleware"
	"strings"
	"time"
)

func InitWebServer(redisClient rc.Cmdable, handler *web.UserHandler) *gin.Engine {
	server := gin.Default()
	middlewares := initMiddlewares(redisClient)
	server.Use(middlewares...)
	handler.RegisterRoutes(server)
	return server
}

func initMiddlewares(redisClient rc.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		corsMiddleware(),
		sessionMiddleware(),
		jwtMiddleware(),
		// ratelimit.NewBuilder(redisClient, time.Second, 100).Build(),
	}
}

func jwtMiddleware() gin.HandlerFunc {
	return middleware.NewLoginJWTMiddlewareBuilder().
		Ignore("/users/signup").
		Ignore("/users/login").
		Ignore("/users/login_sms/code/send").
		Ignore("/users/login_sms").
		Ignore("/hello").
		Build()
}

func corsMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
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
	})
}

func sessionMiddleware() gin.HandlerFunc {
	store, err := redis.NewStore(16, "tcp", config.Config.Redis.Add, "",
		[]byte("MxQP9pSI6BzUL9XVSZrdSeJm6Jbhw42z"), []byte("0XCRv2ip2KMbnId8hT8UowhPc6yiTrhK"))
	if err != nil {
		panic(err)
	}
	// 设置session的存储方式
	return sessions.Sessions("mysessions", store)
}
