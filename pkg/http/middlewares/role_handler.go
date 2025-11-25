package middlewares

import (
	"durich-be/internal/domain"
	"durich-be/pkg/authentication"
	"durich-be/pkg/errors"
	"durich-be/pkg/http/response"

	"github.com/gin-gonic/gin"
)

func RoleHandler(roles ...domain.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		var isEligible bool
		userAuth := authentication.GetUserDataFromToken(c)

		for _, role := range roles {
			for _, userRole := range userAuth.Role {
				if userRole == role {
					isEligible = true
					break
				}
			}
			if isEligible {
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
