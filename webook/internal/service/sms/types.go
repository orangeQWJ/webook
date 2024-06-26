package sms

import "context"

type Service interface {
	//Send(ctx context.Context, numbers []string, appId string, signature string, tpl string, args []string) error
	Send(ctx context.Context, tpl string, args []string, number ...string) error
}
