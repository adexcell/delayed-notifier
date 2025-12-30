package getbyid

import (
	"context"
	"errors"
	"time"

	"github.com/adexcell/delayed-notifier/internal/domain"
)

type Postgres interface {
	GetByID(ctx context.Context, id int64) (*domain.User, error)
}

type Redis interface {
	GetByID(ctx context.Context, userID int64) (*domain.User, error)
	Set(ctx context.Context, user *domain.User, ttl time.Duration) error
}

type Usecase struct {
	postgres Postgres
	redis    Redis
	cacheTTL time.Duration
}

func New(postgres Postgres, redis Redis, cacheTTL time.Duration) *Usecase {
	return &Usecase{
		postgres: postgres,
		redis:    redis,
		cacheTTL: cacheTTL,
	}
}

func (s *Usecase) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	user := &domain.User{}
	user, err := s.redis.GetByID(ctx, id)

	// Cache hit
	if err == nil {
		return user, nil
	}

	if !errors.Is(err, domain.ErrUserNotFound) {
		return nil, domain.ErrCacheFailed
	}

	user, err = s.postgres.GetByID(ctx, id)
	if errors.Is(err, domain.ErrUserNotFound) {
		return nil, err
	}

	if err != nil {
		return nil, domain.ErrInternal
	}

	if err := s.redis.Set(ctx, user, s.cacheTTL); err != nil {
		return nil, domain.ErrCacheFailed
	}

	return user, nil
}
