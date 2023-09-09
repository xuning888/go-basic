package sms

import "context"

// Service 当需要适配多种短信发送服务方时， 应该时刻注意代码重构，防止代码腐化
type Service interface {
	Send(ctx context.Context, tpl string, numbers []string, args []string) error
	SendV1(ctx context.Context, tpl string, numbers []string, args []NamedArg) error
}

// NamedArg 适配阿里云短信发送服务时的新增的参数
type NamedArg struct {
	Val  string
	Name string
}
