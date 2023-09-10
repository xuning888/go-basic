package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go-basic/webbook/internal/domain"
	"time"
)

var ErrKeyNotExit = redis.Nil

type UserCache interface {
	Get(ctx context.Context, id int64) (domain.User, error)
	Set(ctx context.Context, user domain.User) error
}

type redisUserCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

func NewRedisUserCache(client redis.Cmdable) UserCache {
	return &redisUserCache{
		client:     client,
		expiration: time.Minute * 15,
	}
}

func (cache *redisUserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := cache.key(id)
	value, err := cache.client.Get(ctx, key).Bytes()
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	return u, json.Unmarshal(value, &u)
}

func (cache *redisUserCache) Set(ctx context.Context, user domain.User) error {
	jsonStr, err := json.Marshal(user)
	if err != nil {
		return err
	}
	key := cache.key(user.Id)
	return cache.client.Set(ctx, key, jsonStr, cache.expiration).Err()
}

func (cache *redisUserCache) key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}
