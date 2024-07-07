package auth

import (
	"context"
	"errors"
	"xws/webook/internal/service/sms"

	"github.com/golang-jwt/jwt/v5"
)

type SmsService struct {
	svc sms.Service
	key string
}

// send发送, 其中biz必须是线下申请的,代表业务方的token
func (s *SmsService) Send(ctx context.Context, biz string, args []string, number ...string) error {
	// 权限校验
	var tc Claims
	token, err := jwt.ParseWithClaims(biz, &tc, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.key), nil
	})
	if err != nil {
		return errors.New("解析token出错")
	}
	if !token.Valid {
		return errors.New("token 不合法")
	}

	tlp := tc.Tpl
	return s.svc.Send(ctx, tlp, args, number...)

}

type Claims struct {
	jwt.RegisteredClaims
	Tpl string
}

/*
1.提高可用性：重武机制、客戸端限流、failovec （轮询，实时检测）
	1.1 实时检测：
	1.1.1 基于超时的实时检测（连续超时）
	1.1.2 基于应时间的实时检测（比如说，平均响应时间上升 20%）
	1.1.3 基于长尾请求的实时检测（比如说，响应时间超过 1S 的请求占比超过了 10％）
2. 提高安全性:
	2.1 完整的资源申请与审批流程
	2.2 鉴权
	2.2.1 静态token
3. 提高可观测行: 日志, metrics, tracing, 丰富完善的排查手段
4. 可测试性
*/
