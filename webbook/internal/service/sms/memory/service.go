package memory

import (
	"context"
	"go-basic/webbook/internal/service/sms"
	"log"
)

var _ sms.Service = &Service{}

// Service 用于测试的短信发送, 就是简单的在控制台打印一下数据
type Service struct {
}

func NewMemService() sms.Service {
	return &Service{}
}

func (s *Service) SendV1(ctx context.Context, tpl string, numbers []string, args []sms.NamedArg) error {
	log.Printf("numbers: %s, msg: %s", numbers, args)
	return nil
}

func (s *Service) Send(ctx context.Context, tpl string, numbers []string, args []string) error {
	log.Printf("numbers: %s, msg: %s", numbers, args)
	return nil
}
