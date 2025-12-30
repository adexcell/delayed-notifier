package register

import (
	"errors"
	"net/http"

	"github.com/adexcell/delayed-notifier/internal/domain"
	"github.com/gin-gonic/gin"
)

var usecase *Usecase

func HTTPv1(c *gin.Context) {
	input := Input{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := usecase.Register(c, input)
	if errors.Is(err, domain.ErrEmailAlreadyRegistered) {
		c.JSON(http.StatusConflict, nil)
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSONP(http.StatusCreated, output)
}
