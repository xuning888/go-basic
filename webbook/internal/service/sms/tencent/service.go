package tencent

import (
	"context"
	"fmt"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	mysms "go-basic/webbook/internal/service/sms"
	"go-basic/webbook/internal/util"
)

type Service struct {
	appId     *string
	signature *string
	client    *sms.Client
}

func NewService() mysms.Service {
	return nil
}

func (s *Service) Send(ctx context.Context, tpl string, numbers []string, args []string) error {
	req := sms.NewSendSmsRequest()
	req.SmsSdkAppId = s.appId
	req.SignName = s.signature
	req.TemplateId = util.ToPtr[string](tpl)
	req.PhoneNumberSet = s.toStringPtrSlice(numbers)
	req.TemplateParamSet = util.Map[string, *string](args, func(idx int, src string) *string {
		return util.ToPtr[string](src)
	})
	resp, err := s.client.SendSms(req)
	if err != nil {
		return err
	}

	// TODO, 批量发送之后返回批量错误 暂时不用管
	for _, status := range resp.Response.SendStatusSet {
		if status.Code == nil || *(status.Code) != "Ok" {
			return fmt.Errorf("发送短信失败 %s, %s", *status.Code, *status.Message)
		}
	}
	return nil
}

func (s *Service) toStringPtrSlice(src []string) []*string {
	return util.Map[string, *string](src, func(idx int, src string) *string {
		return util.ToPtr[string](src)
	})
}
