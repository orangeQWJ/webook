package ratelimit

import (
	"context"
	"fmt"
	"xws/webook/internal/service/sms"
	"xws/webook/pkg/ratelimit"
)

type RateLimitdSmsService struct {
	svc     sms.Service
	limiter ratelimit.Limiter
}

var errLimited = fmt.Errorf("短信触发限流")

var _ sms.Service = &RateLimitdSmsService{}

func NewRateLimitdSmsService(svc sms.Service, limiter ratelimit.Limiter) sms.Service {
	return &RateLimitdSmsService{
		svc:     svc,
		limiter: limiter,
	}
}

func (s *RateLimitdSmsService) Send(ctx context.Context, tpl string, args []string, number ...string) error {
	// 在这里加一些代码
	for _, phoneNum := range number {
		ok, err := s.limiter.Limit(ctx, phoneNum)
		if err != nil {
			// redis崩了
			// 可以限流: 保守策略, 当下游代码比较弱,可能无法承受大量请求
			// 不限流: 下游很强,能承担. 或者 业务可用性要求很高, 尽量容错策略
			return fmt.Errorf("短信服务判断是否限流出现问题%w", err)
		}
		if ok {
			return errLimited
		}
	}

	err := s.svc.Send(ctx, tpl, args, number...)
	// 在这里加一些代码,新特性
	return err
}
