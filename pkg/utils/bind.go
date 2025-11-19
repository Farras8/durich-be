package utils

import (
	"github.com/gin-gonic/gin"
)

func BindData(ctx *gin.Context, data interface{}) error {
	if err := ctx.ShouldBindJSON(data); err != nil {
		return err
	}
	return nil
}
