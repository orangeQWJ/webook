package failover

import (
	"context"
	"errors"
	"log"
	"sync/atomic"
	"xws/webook/internal/service/sms"
)

type FailoverSmsService struct {
	scvs []sms.Service
	idx  uint64
}

func NewFailoverSmsService(scvs []sms.Service) sms.Service {
	return &FailoverSmsService{
		scvs: scvs,
	}

}

func (f *FailoverSmsService) Send(ctx context.Context, tpl string, args []string, number ...string) error {
	for _, svc := range f.scvs {
		err := svc.Send(ctx, tpl, args, number...)
		// 发送成功
		if err == nil {
			return nil
		}
		// 输出日志
		// 做好监控
		log.Println(err)
	}
	// 网络崩了?
	return errors.New("全部服务商都失败了")
}
func (f *FailoverSmsService) Sendv1(ctx context.Context, tpl string, args []string, number ...string) error {
	idx := atomic.AddUint64(&f.idx, 1)
	length := uint64(len(f.scvs))
	for i := idx; i < idx+length; i++ {
		svc := f.scvs[i%length]
		err := svc.Send(ctx, tpl, args, number...)
		switch err {
		case nil:
			return nil
		case context.DeadlineExceeded, context.Canceled:
			return err
		}
		// 其他情况打印日志
	}
	return errors.New("全部服务商都失败了")
}
