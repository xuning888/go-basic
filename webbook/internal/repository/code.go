package repository

import (
	"context"
	"go-basic/webbook/internal/repository/cache"
)

var _ CodeRepository = &CacheCodeRepository{}

type CodeRepository interface {
	Store(ctx context.Context, biz string, phone string, code string) error
	Verify(ctx context.Context, biz string, phone string, code string) (bool, error)
}

type CacheCodeRepository struct {
	cache cache.CodeCache
}

func NewCodeRepository(codeCache cache.CodeCache) CodeRepository {
	return &CacheCodeRepository{
		cache: codeCache,
	}
}

func (c *CacheCodeRepository) Store(ctx context.Context, biz string, phone string, code string) error {
	return c.cache.Set(ctx, biz, phone, code)
}

func (c *CacheCodeRepository) Verify(ctx context.Context, biz string, phone string, code string) (bool, error) {
	return c.cache.Verify(ctx, biz, phone, code)
}
