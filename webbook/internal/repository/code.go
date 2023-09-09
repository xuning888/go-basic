package repository

import (
	"context"
	"go-basic/webbook/internal/repository/cache"
)

type CodeRepository struct {
	cache cache.CodeCache
}

func NewCodeRepository(codeCache cache.CodeCache) *CodeRepository {
	return &CodeRepository{
		cache: codeCache,
	}
}

func (c *CodeRepository) Store(ctx context.Context, biz string, phone string, code string) error {
	return c.cache.Set(ctx, biz, phone, code)
}

func (c *CodeRepository) Verify(ctx context.Context, biz string, phone string, code string) (bool, error) {
	return c.cache.Verify(ctx, biz, phone, code)
}
