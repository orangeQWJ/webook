package web

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

type JwtHandler struct {
	// access_token key
	AtKey []byte
	// refresh_token key
	RtKey []byte
}

func NewJwtHandler() *JwtHandler {
	return &JwtHandler{
		AtKey: []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"),
		RtKey: []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"),
	}
}

func (j JwtHandler) SetRefreshToken(ctx *gin.Context, userId int64) error {
	claims := RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 3000)),
		},
		Uid: userId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(j.RtKey)
	if err != nil {
		return err
	}
	ctx.Header("x-refresh-token", tokenStr)
	return nil
}

func (j JwtHandler) SetJWT(ctx *gin.Context, userId int64) error {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		Uid:       userId,
		UserAgent: ctx.Request.UserAgent(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(j.AtKey)
	if err != nil {
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

func ExtractToken(ctx *gin.Context) string {
	tokenHeader := ctx.GetHeader("Authorization")
	if tokenHeader == "" {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return ""
	}
	segs := strings.Split(tokenHeader, " ")
	if len(segs) != 2 {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return ""
	}
	tokenStr := segs[1]
	return tokenStr
}

type RefreshClaims struct {
	Uid int64
	jwt.RegisteredClaims
}

type UserClaims struct {
	jwt.RegisteredClaims
	// 声明你自己要放进token里的数据
	Uid int64
	// 自己随便加
	// 敏感信息不要加
	UserAgent string
}
