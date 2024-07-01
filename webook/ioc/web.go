package ioc

import (
	"strings"
	"time"
	"xws/webook/internal/web"
	"xws/webook/internal/web/middleware"
	"xws/webook/pkg/ginx/middlewares/ratelimit"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func InitWebServer(mdls []gin.HandlerFunc, hdl *web.UserHandler) *gin.Engine {
	service := gin.Default()
	service.Use(mdls...)
	hdl.RegisterRoutes(service)
	return service
}

func InitMiddlewares(redisClient redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		corsHdl(),
		jwtHdl(),
		ratelimitHdl(redisClient),
	}
}

func ratelimitHdl(redisClient redis.Cmdable) gin.HandlerFunc {
	return ratelimit.NewBuilder(redisClient, time.Minute, 100).Build()
}

func jwtHdl() gin.HandlerFunc {
	return middleware.NewLoginJwtMiddlewareBuilder().
		IgnorePaths("/users/signup").
		IgnorePaths("/users/login").
		IgnorePaths("/users/login_sms/code/send").
		IgnorePaths("/users/login_sms").
		Build()
}

func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		//AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"POST", "GET"},
		// 允许前端发送的字段
		AllowHeaders: []string{"authorization", "content-type"},
		//ExposeHeaders:    []string{"authorization", "content-type"},
		//authorization,content-type
		// 不加这个,前端拿不到
		ExposeHeaders:    []string{"x-jwt-token"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			//return origin == "https://github.com"
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "webook.com")
		},
		MaxAge: 12 * time.Hour,
	})
}
