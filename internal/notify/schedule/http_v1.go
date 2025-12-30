package schedule

import (
	"net/http"

	"github.com/adexcell/delayed-notifier/internal/domain"
	"github.com/gin-gonic/gin"
)

var usecase *Usecase

func HTTPv1(c *gin.Context) {
	val, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	userID, ok := val.(int64)
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	input := &Input{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	n := &domain.Notify{
		UserID:  userID,
		Message: input.Message,
		SendAt:  input.SendAt,
	}

	if err := usecase.Schedule(c, n); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

}
