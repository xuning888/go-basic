package sms

import "context"

type Service interface {
	Send(ctx context.Context, tpl string, numbers []string, args []string) error
}

type NamedArg struct {
	Val  string
	Name string
}
