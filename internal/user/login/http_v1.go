package login

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var usecase *Usecase

func HTTPv1(c *gin.Context) {
	input := Input{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := usecase.Login(c, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSONP(http.StatusOK, output)
}
