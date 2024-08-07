package tencent

import (
	"context"
	"fmt"
	"os"

	"github.com/ecodeclub/ekit/slice"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	Tencent_sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Service struct {
	appId    *string
	signName *string
	client   *Tencent_sms.Client
	//limiter  ratelimit.Limiter
}

func ToPtr(t string) *string {
	return &t
}

/*
func NewService(client *Tencent_sms.Client, appId string, signName string) *Service {
	return &Service{
		appId:    ToPtr(appId),
		signName: ToPtr(signName),
		client:   client,
	}
}
*/

//func NewService(client *Tencent_sms.Client, limiter ratelimit.Limiter) *Service {
func NewService(client *Tencent_sms.Client) *Service {
	return &Service{
		appId:    ToPtr("1400920455"),
		signName: ToPtr("木凳也公众号"),
		client:   client,
		//limiter:  limiter,
	}
}
func InitTencentSmsClient() *Tencent_sms.Client {
	secretId, ok := os.LookupEnv("SMS_SECRET_ID")
	if !ok {
		panic("设置SMS_SECRET_ID")
	}
	secretKey, ok := os.LookupEnv("SMS_SECRET_KEY")
	if !ok {
		panic("设置SMS_SECRET_KEY")
	}
	c, err := Tencent_sms.NewClient(common.NewCredential(secretId, secretKey),
		"ap-beijing",
		profile.NewClientProfile())
	if err != nil {
		panic("腾讯短信服务客户端启动失败")
	}
	return c
}

func (s *Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	// 腾讯云在向某手机号发送短信时有频率限制,在网页端可以修改限制
	// 将某些手机号加入无限制的白名单,方便测试开发
	req := Tencent_sms.NewSendSmsRequest()
	req.SmsSdkAppId = s.appId
	req.SignName = s.signName
	req.TemplateId = ToPtr(tplId)
	req.PhoneNumberSet = s.toStringPtrSlice(numbers)
	req.TemplateParamSet = s.toStringPtrSlice(args)
	resp, err := s.client.SendSms(req)
	if err != nil {
		return err
	}
	for _, status := range resp.Response.SendStatusSet {
		if status.Code == nil || *(status.Code) != "Ok" {
			return fmt.Errorf("发送短信失败 %s, %s", *status.Code, *status.Message)
		}
	}
	return nil
}

func (s *Service) toStringPtrSlice(src []string) []*string {
	return slice.Map[string, *string](src, func(idx int, src string) *string {
		return &src
	})
}
