package jwt

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	AtKey = []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0")
	RtKey = []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0")
)

type RedisJWTHandler struct {
	Cmd redis.Cmdable
}

var _ Handler = &RedisJWTHandler{}

func NewRedisJWTHandler(cmd redis.Cmdable) Handler {
	return &RedisJWTHandler{
		Cmd: cmd,
	}
}

// 设置短token,nil 设置成功,否则设置失败
func (h *RedisJWTHandler) SetJWT(ctx *gin.Context, userId int64, ssid string) error {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		Uid:       userId,
		Ssid:      ssid,
		UserAgent: ctx.Request.UserAgent(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(AtKey)
	if err != nil {
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

// 设置长token,nil 设置成功,否则设置失败
func (h *RedisJWTHandler) SetRefreshToken(ctx *gin.Context, userId int64, ssid string) error {
	claims := RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
		Ssid: ssid,
		Uid:  userId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(RtKey)
	if err != nil {
		return err
	}
	ctx.Header("x-refresh-token", tokenStr)
	return nil
}

// 同时设置长短token, 返回nil 设置成功,返回其他设置失败
func (h *RedisJWTHandler) SetLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := h.SetJWT(ctx, uid, ssid)
	if err != nil {
		return err
	}
	err = h.SetRefreshToken(ctx, uid, ssid)
	return err

}

// 提取Http请求的 Authorization 字段
func (h *RedisJWTHandler) ExtractToken(ctx *gin.Context) string {
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

// 通过将token中的ssid字段存于redis来实现退出登录.redis中记录已经退出登录的ssid(标识失效token)
// nil 退出登录成功, 否则失败
func (h *RedisJWTHandler) ClearToken(ctx *gin.Context) error {
	ctx.Header("x-jwt-token", "")
	ctx.Header("x-refresh-token", "")
	refreshToken := h.ExtractToken(ctx) // 约定在Authorization字段存放refreshToken
	var reFreshClaims RefreshClaims
	token, err := jwt.ParseWithClaims(refreshToken, &reFreshClaims, func(t *jwt.Token) (interface{}, error) {
		return RtKey, nil
	})

	if err != nil || !token.Valid {
		return err
	}

	remainTime := reFreshClaims.RegisteredClaims.ExpiresAt.Sub(time.Now())
	err = h.Cmd.Set(ctx, fmt.Sprintf("users:ssid:%s", reFreshClaims.Ssid), "", remainTime).Err()
	return err
}
func (j *RedisJWTHandler) CheckSession(ctx *gin.Context, ssid string) error {
	_, err := j.Cmd.Exists(ctx, fmt.Sprintf("users:ssid:%s", ssid)).Result()
	return err
	/*

		if err != nil || count > 0 {
			// 要么redis有问题,要么已经退出登录了
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	*/

}
