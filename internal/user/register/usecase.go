package register

import (
	"context"
	"errors"
	"time"

	"github.com/adexcell/delayed-notifier/internal/adapter/postgres"
	"github.com/adexcell/delayed-notifier/internal/domain"
	"github.com/adexcell/delayed-notifier/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type Postgres interface {
	CreateUser(ctx context.Context, email string, passwordHash string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
}

type TokenManager interface {
	NewJWT(userID int64, ttl time.Duration) (string, error)
	Parse(accessToken string) (int64, error)
}

type Usecase struct {
	postgres     Postgres
	tokenManager TokenManager
	tokenTTL     time.Duration
}

func New(postgres *postgres.Postgres, tokenManager *auth.TokenManager, tokenTTL time.Duration) *Usecase {
	return &Usecase{
		postgres: postgres,
		tokenManager: *tokenManager,
		tokenTTL: tokenTTL,
	}
}

func (c *Usecase) Register(ctx context.Context, input Input) (Output, error) {
	_, err := c.postgres.GetUserByEmail(ctx, input.Email)
	if err == nil {
		return Output{}, domain.ErrEmailAlreadyRegistered
	}

	if !errors.Is(err, domain.ErrUserNotFound) {
		return Output{}, domain.ErrInternal
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if errors.Is(err, bcrypt.ErrPasswordTooLong) {
		return Output{}, domain.ErrPasswordTooLong
	} else if err != nil {
		return Output{}, domain.ErrInternal
	}

	user, err := c.postgres.CreateUser(ctx, input.Email, string(passwordHash))
	if errors.Is(err, domain.ErrInternal) {
		return Output{}, err
	}

	token, err := c.tokenManager.NewJWT(user.ID, c.tokenTTL)
	if err != nil {
		return Output{}, domain.ErrInvalidCredentials
	}

	return Output{
		User: user,
		AccessToken: token,
	}, nil
}
