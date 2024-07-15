package middleware

import (
	"net/http"

	ijwt "xws/webook/internal/web/jwt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWT 登录校验
type LoginJwtMiddlewareBuilder struct {
	paths []string
	ijwt.Handler
}

func NewLoginJwtMiddlewareBuilder(jwtHdl ijwt.Handler) *LoginJwtMiddlewareBuilder {
	return &LoginJwtMiddlewareBuilder{
		Handler: jwtHdl,
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
		tokenStr := l.ExtractToken(ctx)
		claims := ijwt.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (interface{}, error) {
			return ijwt.AtKey, nil
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
		// 如果redis崩溃了,可以采用降级服务
		// 不影响正常用户的登录
		// 正常用户退出登录时,我会清掉他的header中的token
		// 他不会走到这一步
		// 黑客如果把redis搞崩,就可以用已过期的token来登录
		err = l.CheckSession(ctx, claims.Ssid)
		if err != nil {
			// 要么redis有问题,要么已经退出登录了
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
