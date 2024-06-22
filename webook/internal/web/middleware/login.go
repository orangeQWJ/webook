package middleware

import (
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}

}

func (l *LoginMiddlewareBuilder) IgnorePaths(path string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		/*
			if ctx.Request.URL.Path == "/users/login" || ctx.Request.URL.Path == "/users/signup" {
				return
			}
		*/
		sess := sessions.Default(ctx)
		if id := sess.Get("userId"); id == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		updateTime := sess.Get("update_time")
		now := time.Now().UnixMilli()
		// 刚登陆,还没刷新过
		if updateTime == nil {
			sess.Set("update_time", now)
			sess.Options(sessions.Options{
				MaxAge: 60 * 5,
			})
			sess.Save()
			return
		}
		// 已经更新过了
		updateTimeVal, ok := updateTime.(int64)
		if !ok {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if now-updateTimeVal > 60*1000 {
			sess.Set("update_time", now)
			sess.Options(sessions.Options{
				MaxAge: 60 * 5,
			})
			sess.Save()
			return
		}
	}
}
