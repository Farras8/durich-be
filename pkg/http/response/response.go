package response

import (
	"net/http"

	"durich-be/pkg/errors"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Code    int         `json:"code"`
}

func SendSuccess(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, Response{
		Success: true,
		Message: message,
		Data:    data,
		Code:    code,
	})
}

func SendError(c *gin.Context, err error) {
	if appErr, ok := err.(*errors.AppError); ok {
		c.JSON(appErr.Code, Response{
			Success: false,
			Message: appErr.Message,
			Code:    appErr.Code,
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusInternalServerError, Response{
		Success: false,
		Message: err.Error(),
		Code:    http.StatusInternalServerError,
	})
	c.Abort()
}
