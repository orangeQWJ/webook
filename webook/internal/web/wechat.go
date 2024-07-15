package web

import (
	"net/http"
	"xws/webook/internal/service"
	"xws/webook/internal/service/oauth2/wechat"
	ijwt "xws/webook/internal/web/jwt"

	"github.com/gin-gonic/gin"
)

type OAuth2WechatHandler struct {
	svc     wechat.Service
	userSvc service.UserService
	ijwt.Handler
}

func NewOAuth2WechatHandler(svc wechat.Service, userSvc service.UserService, jwtHdl ijwt.Handler) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc: svc,
		userSvc: userSvc,
		Handler: jwtHdl,
	}
}

func (h *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/oauth2/wechat")
	g.GET("/authurl", h.AuthURL)
	g.Any("/callback", h.Callback)

}

func (h *OAuth2WechatHandler) AuthURL(ctx *gin.Context) {
	url, err := h.svc.AuthURL(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "构造扫码登录URL失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: url,
	})

}
func (h *OAuth2WechatHandler) Callback(ctx *gin.Context) {
	// 验证微信的code
	code := ctx.Query("code")
	state := ctx.Query("state")
	wechatInfo, err := h.svc.VerifyCode(ctx, code, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	u, err := h.userSvc.FindOrCreateByWechat(ctx, wechatInfo)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	err = h.SetLoginToken(ctx, u.Id)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "JWT系统错误")
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})
}
