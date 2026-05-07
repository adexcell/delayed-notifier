package dto

// DeleteNotifyInput represents the input data required to delete a notification.
type DeleteNotifyInput struct {
	ID string `binding:"required" json:"id"`
}
