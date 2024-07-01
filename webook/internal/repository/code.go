package repository

import (
	"context"
	"xws/webook/internal/repository/cache"
)

var (
	ErrCodeSendTooMany = cache.ErrCodeSendTooMany
	ErrCodeExpired     = cache.ErrCodeExpired
)

var _ CodeRepository = (*CacheCodeRepository)(nil)

type CodeRepository interface {
	Store(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

type CacheCodeRepository struct {
	cache cache.CodeCache
}

func NewCodeRepository(c cache.CodeCache) CodeRepository {
	return &CacheCodeRepository{
		cache: c,
	}
}

func (repo *CacheCodeRepository) Store(ctx context.Context, biz, phone, code string) error {
	return repo.cache.Set(ctx, biz, phone, code)
}

func (repo *CacheCodeRepository) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	return repo.cache.Verify(ctx, biz, phone, inputCode)
}
