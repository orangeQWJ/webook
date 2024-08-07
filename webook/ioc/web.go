package ioc

import (
	"strings"
	"time"
	"xws/webook/internal/web"
	ijwt "xws/webook/internal/web/jwt"
	"xws/webook/internal/web/middleware"
	"xws/webook/pkg/ginx/middlewares/ratelimit"
	limit "xws/webook/pkg/ratelimit"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func InitWebServer(mdls []gin.HandlerFunc, userHdl *web.UserHandler, oauth2WechatHdl *web.OAuth2WechatHandler) *gin.Engine {
	service := gin.Default()
	service.Use(mdls...)
	userHdl.RegisterRoutes(service)
	oauth2WechatHdl.RegisterRoutes(service)
	return service
}

func InitMiddlewares(redisClient redis.Cmdable, jwtHdl ijwt.Handler) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		corsHdl(),
		GenJwtHdl(jwtHdl),
		ratelimitHdl(redisClient),
	}
}

func ratelimitHdl(redisClient redis.Cmdable) gin.HandlerFunc {
	//return ratelimit.NewBuilder(redisClient, time.Minute, 100).Build()
	return ratelimit.NewBuilder(limit.NewRedisSlidingWindowLimiter(redisClient, time.Minute, 100)).Build()
}

func GenJwtHdl(jwtHdl ijwt.Handler) gin.HandlerFunc {
	return middleware.NewLoginJwtMiddlewareBuilder(jwtHdl).
		IgnorePaths("/users/signup").
		IgnorePaths("/users/login").
		IgnorePaths("/users/login_sms/code/send").
		IgnorePaths("/users/login_sms").
		IgnorePaths("/oauth2/wechat/authurl").
		IgnorePaths("/oauth2/wechat/callback").
		IgnorePaths("/users/refresh_token").
		IgnorePaths("/users/LogoutJWT").
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
		ExposeHeaders:    []string{"x-jwt-token", "x-refresh_token"},
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
