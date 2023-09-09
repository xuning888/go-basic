package service

import (
	"context"
	"fmt"
	"go-basic/webbook/internal/repository"
	"go-basic/webbook/internal/service/sms"
	"math/rand"
)

var tpl = "0"

type CodeService struct {
	repo   *repository.CodeRepository
	smsSvc sms.Service
}

func NewCodeService(codeRepository *repository.CodeRepository, smsSvc sms.Service) *CodeService {
	return &CodeService{
		repo:   codeRepository,
		smsSvc: smsSvc,
	}
}

func (c *CodeService) Send(ctx context.Context, biz string, phone string) error {

	code := c.code()
	err := c.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}

	err = c.smsSvc.Send(ctx, tpl, []string{phone}, []string{code})
	// if err != nil {
	// 意味着，redis 有个验证码
	// 我们能不能删除这个验证码?
	// 这个err 可能是超时的err, 我们无法明确发出去了没, 所以redis中的key不能删除
	// }
	return err
}

func (c *CodeService) code() string {
	// 生成 0 ~ 999999的随机数
	num := rand.Intn(1000000)
	return fmt.Sprintf("%6d", num)
}
