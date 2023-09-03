package memory

import (
	"context"
	"log"
)

type Service struct {
}

func (s *Service) Send(ctx context.Context, tpl string, numbers []string, args []string) error {
	log.Printf("msg: %s\n", args)
	return nil
}
