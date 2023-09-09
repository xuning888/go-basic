package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
)

var (
	ErrSendCodeToMany         = errors.New("发送验证码太频繁")
	ErrCodeVerifyTooManyTimes = errors.New("验证次数太多了")
	ErrUnknownForCode         = errors.New("我也不知发生什么了，反正是跟 code 有关")
)

//go:embed lua/send_code.lua
var luaCodeSet string

//go:embed lua/verify_code.lua
var luaVerifySet string

type CodeCache interface {
	Set(ctx context.Context, biz string, phone string, code string) error
	Verify(ctx context.Context, biz string, phone string, code string) (bool, error)
}

type RedisCodeCache struct {
	client redis.Cmdable
}

func NewRedisCodeCache(client redis.Cmdable) CodeCache {
	return &RedisCodeCache{
		client: client,
	}
}

func (c *RedisCodeCache) Set(ctx context.Context, biz string, phone string, code string) error {
	key := c.key(biz, phone)
	log.Printf("sendKey: %s\n", key)
	res, err := c.client.Eval(ctx, luaCodeSet, []string{key}, code).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		return nil
	case -1:
		return ErrSendCodeToMany
	default:
		return errors.New("系统错误")
	}
}

func (c *RedisCodeCache) Verify(ctx context.Context, biz string, phone string, code string) (bool, error) {
	key := c.key(biz, phone)
	log.Printf("verifyKey: %s\n", key)
	res, err := c.client.Eval(ctx, luaVerifySet, []string{key}, code).Int()
	if err != nil {
		return false, err
	}

	switch res {
	case 0:
		// 验证通过
		return true, nil
	case -1:
		// 改验证码已经验证过了， 或者改验证码已经过期了
		return false, ErrCodeVerifyTooManyTimes
	case -2:
		return false, nil
	}

	return false, ErrUnknownForCode
}

func (c *RedisCodeCache) key(biz string, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
