package schedule

import "time"

type Output struct {
	
}

type Input struct {
	Message string    `json:"message" binding:"required"`
	SendAt  time.Time `json:"send_at" binding:"required"`
}
