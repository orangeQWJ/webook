package middleware

import (
	"net/http"

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

	}
}
