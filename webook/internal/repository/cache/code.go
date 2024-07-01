package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var (
	ErrCodeSendTooMany    = errors.New("发送验证码太频繁")
	ErrCodeExpired        = errors.New("验证码失效")
	ErrUnknowForLuaScript = errors.New("lua 脚本出现了未知的错误")
	ErrEnryWithoutExpire  = errors.New("验证码数据没有过期时间")
)

// 编译器会把lua脚本的源码放到这个变量中
//
//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/verify_code.lua
var luaVerifyCode string

var _ CodeCache = &RedisCodeCache{}

type CodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

type RedisCodeCache struct {
	client redis.Cmdable
}

func NewCodeCache(client redis.Cmdable) CodeCache {
	return &RedisCodeCache{
		client: client,
	}
}

func (c *RedisCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	res, err := c.client.Eval(ctx, luaSetCode, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		return nil
	case -1:
		// 发送太频繁
		return ErrCodeSendTooMany
	case -2:
		//系统错误
		return ErrEnryWithoutExpire
	default:
		//lua脚本执行出错系统错误
		return ErrUnknowForLuaScript
	}
}

func (c *RedisCodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

func (c *RedisCodeCache) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	res, err := c.client.Eval(ctx, luaVerifyCode, []string{c.key(biz, phone)}, inputCode).Int()
	if err != nil {
		return false, err
	}
	switch res {
	case 0:
		return true, nil
	case -1:
		//
		return false, ErrCodeExpired
	case -2:
		return false, nil
	default:
		return false, ErrUnknowForLuaScript
	}
}
