package login

import (
	"context"
	"errors"
	"time"

	"github.com/adexcell/delayed-notifier/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type Postgres interface {
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

type Redis interface {
	GetByID(ctx context.Context, userID int64) (*domain.User, error)
	Set(ctx context.Context, user *domain.User, ttl time.Duration) error
}

type TokenManager interface {
	NewJWT(userID int64, ttl time.Duration) (string, error)
	Parse(accessToken string) (int64, error)
}

type Usecase struct {
	postgres     Postgres
	redis        Redis
	tokenManager TokenManager
	tokenTTL     time.Duration
}

func New(postgres Postgres, redis Redis) *Usecase {
	return &Usecase{
		postgres: postgres,
		redis:    redis,
	}
}

func (c *Usecase) Login(ctx context.Context, input Input) (Output, error) {
	user, err := c.postgres.GetByEmail(ctx, input.Email)

	if errors.Is(err, domain.ErrUserNotFound) {
		return Output{}, domain.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return Output{}, domain.ErrInvalidCredentials
	}

	token, err := c.tokenManager.NewJWT(user.ID, c.tokenTTL)
	if err != nil {
		return Output{}, domain.ErrInvalidCredentials
	}

	return Output{
		Token: token,
	}, nil
}
