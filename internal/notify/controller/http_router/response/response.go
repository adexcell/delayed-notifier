package response

import (
	"net/http"

	"github.com/wb-go/wbf/ginext"
)

type Response struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"error,omitempty"`
	Data         any    `json:"data,omitempty"`
}

func Created(c *ginext.Context, data any) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    data,
	})
}

func OK(c *ginext.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

func NoContent(c *ginext.Context) {
	c.JSON(http.StatusBadRequest, Response{
		Success: true,
	})
}

func ValidationsError(c *ginext.Context, data any) {
	c.AbortWithStatusJSON(http.StatusBadRequest, Response{
		Success:      false,
		ErrorMessage: "validation error",
		Data:         data,
	})
}

func BadRequest(c *ginext.Context, data any) {
	c.AbortWithStatusJSON(http.StatusBadRequest, Response{
		Success:      false,
		ErrorMessage: "request failed",
		Data:         data,
	})
}

func InternalServerError(c *ginext.Context, data any) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, Response{
		Success:      false,
		ErrorMessage: "request failed",
		Data:         data,
	})
}
