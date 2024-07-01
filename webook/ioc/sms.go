package ioc

import (
	"xws/webook/internal/service/sms"
	"xws/webook/internal/service/sms/tencent"

	Tencent_sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

func InitSMSService(client *Tencent_sms.Client) sms.Service {
	return tencent.NewService(client)
}
