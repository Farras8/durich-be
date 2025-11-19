package middlewares

import (
	"durich-be/pkg/authentication"
	"durich-be/pkg/errors"
	"durich-be/pkg/http/response"

	"github.com/gin-gonic/gin"
)

func RoleHandler(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var isEligible bool
		userAuth := authentication.GetUserDataFromToken(c)

		for _, role := range roles {
			if userAuth.Role == role {
				isEligible = true
				break
			}
		}

		if !isEligible {
			response.SendError(c, errors.ForbiddenErrorToAppError())
			c.Abort()
			return
		}

		c.Next()
	}
}
