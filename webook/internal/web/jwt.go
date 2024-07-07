package web

import (
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

type jwtHandler struct{

}

func (j jwtHandler) SetJWT(ctx *gin.Context, userId int64) error {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		Uid:       userId,
		UserAgent: ctx.Request.UserAgent(),
	}

	//token := jwt.New(jwt.SigningMethodHS512)
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"))
	if err != nil {
		//ctx.String(http.StatusInternalServerError, "JWT系统错误")
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

type UserClaims struct {
	jwt.RegisteredClaims
	// 声明你自己要放进token里的数据
	Uid int64
	// 自己随便加
	// 敏感信息不要加
	UserAgent string
}
