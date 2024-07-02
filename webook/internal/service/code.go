package service

import (
	"context"
	"fmt"
	"math/rand"
	"xws/webook/internal/repository"
	"xws/webook/internal/service/sms"
)

const codeTplId = "2196630"

var (
	ErrCodeSendTooMany = repository.ErrCodeSendTooMany
	ErrCodeExpired     = repository.ErrCodeExpired
)

var _ CodeService = &codeService{}

type CodeService interface {
	Send(ctx context.Context, biz string, phone string) error
	Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error)
	GenerateCode() string
}

type codeService struct {
	repo   repository.CodeRepository
	smsSvc sms.Service //接口, tencent.Service 实现了这个接口
}

func NewCodeService(repo repository.CodeRepository, smsSvc sms.Service) CodeService {
	return &codeService{
		repo:   repo,
		smsSvc: smsSvc,
	}
}

func (svc *codeService) Send(ctx context.Context, biz string, phone string) error {
	// biz 区别业务场景
	// phone_code:$biz:130xxxxxx
	// $biz:code:130xxxxxx
	// 1. 生成验证码
	code := svc.GenerateCode()
	// 2. 存入redis
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	// 3. 发送验证码
	err = svc.smsSvc.Send(ctx, codeTplId, []string{code}, phone)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func (svc *codeService) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	return svc.repo.Verify(ctx, biz, phone, inputCode)
}

func (svc *codeService) GenerateCode() string {
	num := rand.Intn(1000000)
	return fmt.Sprintf("%06d", num)
}
