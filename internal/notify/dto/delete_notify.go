package dto

import "github.com/google/uuid"

type DeleteNotifyInput struct {
	ID uuid.UUID `json:"id"`
}

