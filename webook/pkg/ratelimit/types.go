package ratelimit

import "context"

type Limiter interface {
	// Limited 有没有触发限流,key 就是限流对象
	// bool: true 要限流
	// err: 限流器本身有咩有错误
	Limit(ctx context.Context, key string) (bool, error)
}
