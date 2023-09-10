package dal

import (
	"go-basic/webbook/internal/service/sms"
	"go-basic/webbook/internal/service/sms/memory"
)

func InitSmsService() sms.Service {
	return memory.NewMemService()
}
