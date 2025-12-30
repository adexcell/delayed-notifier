package domain

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type User struct {
	ID           int64     `json:"id" db:"id"`
	Email        string    `json:"email" db:"email" validate:"email"`
	TelegramID   int64     `json:"tg_id"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at" `
}

var validate = validator.New(validator.WithPrivateFieldValidation())

func NewUser(ID int64, email string, tgID int64, passwordHash string, createdAt time.Time) (*User, error) {
	u := User{
		ID:           ID,
		Email:        email,
		TelegramID:   tgID,
		PasswordHash: passwordHash,
		CreatedAt:    createdAt,
	}

	if err := u.Validate(); err != nil {
		return &User{}, fmt.Errorf("u.Validate: %w", err)
	}

	return &u, nil
}

func (u *User) Validate() error {
	if err := validate.Struct(u); err != nil {
		return fmt.Errorf("validate.Struct: %w", err)
	}

	return nil
}
