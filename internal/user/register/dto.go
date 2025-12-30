package register

import "github.com/adexcell/delayed-notifier/internal/domain"

type Output struct {
	User        *domain.User
	AccessToken string
}

type Input struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}
