package middlewares

import (
	"strings"

	"durich-be/internal/dto/requests"
	"durich-be/pkg/authentication"
	internal_err "durich-be/pkg/errors"
	"durich-be/pkg/http/response"

	"github.com/gin-gonic/gin"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.SendError(
				c,
				internal_err.AuthError(authentication.AuthErrMalformedToken.Error()),
			)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			response.SendError(
				c,
				internal_err.AuthError(authentication.AuthErrInvalidToken.Error()),
			)
			return
		}

		accessToken, err := authentication.JWTAuth.VerifyAccessToken(tokenString)
		if err != nil {
			response.SendError(c, err)
			return
		}

		c.Set(authentication.Token, requests.UserAuth{
			AuthID:         accessToken.AuthID,
			UserID:         accessToken.UserID,
			Email:          accessToken.Email,
			Role:           accessToken.Role,
			RefreshTokenID: accessToken.RefreshTokenID,
		})

		c.Next()
	}
}
