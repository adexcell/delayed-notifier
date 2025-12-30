package postgres

import (
	"context"
	"errors"

	"github.com/adexcell/delayed-notifier/internal/domain"
	"github.com/adexcell/delayed-notifier/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	db *pgxpool.Pool
	l  *logger.Zerolog
}

func NewPostgres(db *pgxpool.Pool, l  *logger.Zerolog) *Postgres {
	return &Postgres{
		db: db,
		l: l,
	}
}

func (p *Postgres) CreateUser(ctx context.Context, email string, passwordHash string) (*domain.User, error) {
	query := `
		insert into users (email, password_hash)
		values ($1, $2)
		returning id, created_at`

	user := &domain.User{}

	err := p.db.QueryRow(ctx, query, email, passwordHash).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		p.l.Error().
			Err(err).
			Any("user", user).
			Msg("postgres scan failed")
		return nil, domain.ErrInternal
	}

	p.l.Info().
		Any("user", user).
		Msg("postgres scan failed")

	p.l.Info().Any("user", user).Msg("success")
	return user, nil
}

func (p *Postgres) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `select id, email, password_hash, created_at from users where email=$1`

	user := &domain.User{}
	err := p.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		p.l.Info().
			Any("user", user).
			Msg("postgres no rows")
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		p.l.Error().
			Err(err).
			Any("user", user).
			Msg("postgres scan failed")
		return nil, domain.ErrInternal
	}

	p.l.Info().Any("user", user).Msg("success")
	return user, nil
}

func (p *Postgres) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `select id, email, password_hash, created_at from users where id=$1`

	user := &domain.User{}
	err := p.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		p.l.Info().
			Any("user", user).
			Msg("postgres no rows")
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		p.l.Error().
			Err(err).
			Any("user", user).
			Msg("postgres scan failed")
		return nil, domain.ErrInternal
	}

	p.l.Info().Any("user", user).Msg("success")
	return user, nil
}

func (p *Postgres) CreateNotification(ctx context.Context, )
