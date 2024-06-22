package middleware

import (
	"net/http"
	"strings"
	"time"
	"xws/webook/internal/web"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWT 登录校验
type LoginJwtMiddlewareBuilder struct {
	paths []string
}

func NewLoginJwtMiddlewareBuilder() *LoginJwtMiddlewareBuilder {
	return &LoginJwtMiddlewareBuilder{}

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
		tokenHeader := ctx.GetHeader("Authorization")
		if tokenHeader == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		segs := strings.Split(tokenHeader, " ")
		if len(segs) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[1]
		claims := &web.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"), nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		/*
			if claims.ExpiresAt.Time.Before(time.Now()) {
				// 过期了
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			//token.Valid 这里已经能反映出是否过期了
		*/

		if token == nil || !token.Valid || claims.Uid == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		now := time.Now()
		// 每十秒刷新一次
		if claims.ExpiresAt.Sub(now) < time.Second*50 {
			newClaims := web.UserClaims{
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
				},
				Uid: claims.Uid,
			}
			//claims.ExpiresAt =  jwt.NewNumericDate(time.Now().Add(time.Minute)),
			newToken := jwt.NewWithClaims(jwt.SigningMethodHS512, newClaims)
			newTokenStr, newErr := newToken.SignedString([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"))
			if newErr != nil {
				// 日志记录
			}
			//fmt.Println(tokenStr)
			ctx.Header("x-jwt-token", newTokenStr)

		}

		//claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
		//ctx.Set("userId", claims.Uid)
		ctx.Set("claims", claims)

	}
}
