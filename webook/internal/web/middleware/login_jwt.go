package middleware

import (
	"net/http"
	"xws/webook/internal/web"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWT 登录校验
type LoginJwtMiddlewareBuilder struct {
	paths []string
	web.JwtHandler
}

func NewLoginJwtMiddlewareBuilder() *LoginJwtMiddlewareBuilder {
	return &LoginJwtMiddlewareBuilder{
		JwtHandler: *web.NewJwtHandler(),
	}

}

func (l *LoginJwtMiddlewareBuilder) IgnorePaths(path string) *LoginJwtMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}
func (l *LoginJwtMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		tokenStr := web.ExtractToken(ctx)
		claims := web.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (interface{}, error) {
			return l.AtKey, nil
		})

		if err != nil || !token.Valid {

			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if claims.UserAgent != ctx.Request.UserAgent() {
			// 严重安全隐患,需要监控
			// todo 日志监控
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 每十分钟刷新一次
		// 在引入长短token后就不需要定期刷新了
		// 短token由专门的接口刷新
		/*
			if claims.ExpiresAt.Sub(time.now()) < time.Minute*20 {
				l.SetJWT(ctx, claims.Uid)
			}
		*/
		ctx.Set("claims", &claims)
	}
}
