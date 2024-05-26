package authService

import (
	"github.com/gin-gonic/gin"
	. "github.com/martynasd123/golang-scraper/services/auth/constants"
	"net/http"
)

func RequireAuth(service *AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken, err := ctx.Cookie(AccessTokenCookieName)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, "access token required")
			return
		}
		username, err := service.ValidateAccessToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, "access token invalid")
			return
		}
		ctx.Set(UserNameContextKey, username)
		ctx.Next()
	}
}
