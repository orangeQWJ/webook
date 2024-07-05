package failover

import (
	"context"
	"errors"
	"sync/atomic"
	"xws/webook/internal/service/sms"
)

type TimeoutFailoverSmsService struct{
	scvs []sms.Service
	idx int32
	// 连续超时的个数
	cnt int32
	// 连续超时次数超过这个次数,就要切换
	threshold int32 
}

func NewTimeoutFailoverSmsService() sms.Service{
	return &TimeoutFailoverSmsService{}
}

func (t * TimeoutFailoverSmsService) Send(ctx context.Context, tpl string, args []string, number ...string) error {
	idx := atomic.LoadInt32(&t.idx)
	cnt := atomic.LoadInt32(&t.cnt)
	if cnt > t.threshold{
		// 切换
		newIndx := (idx + 1) % int32(len(t.scvs))
		if atomic.CompareAndSwapInt32(&t.idx, idx, newIndx){
			atomic.StoreInt32(&t.cnt, 0)
		}
		// idx = newIndx
		idx = atomic.LoadInt32(&t.idx)
	}
	svc := t.scvs[idx]
	err := svc.Send(ctx, tpl, args, number...)
	switch err {
		case context.DeadlineExceeded:
		atomic.AddInt32(&t.cnt, 1)
	case nil:
		atomic.StoreInt32(&t.cnt, 0)
		return nil
	default:
		// 不知道什么错误
		return err

	}
	return errors.New("全部服务商都失败了")
}
