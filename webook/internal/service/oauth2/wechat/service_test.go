package wechat

import (
	"context"
	"fmt"
	"testing"
)

func Test_service_e2e_VerifyCode(t *testing.T){
	svc := NewService("wx7256bc69ab349c72", "appkey")
	res, err := svc.VerifyCode(context.Background(),"001bKN100hrEQ1sNjO0ZSPHf2bKN1u", "state" )
	fmt.Println(err)
	fmt.Println(res)
}
