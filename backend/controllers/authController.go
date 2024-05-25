package authController

import (
	"net/http"

	"github.com/gin-gonic/gin"
	authService "github.com/martynasd123/golang-scraper/services/auth"
)

type AuthController struct {
	service *authService.AuthService
}

func NewAuthController(service *authService.AuthService) *AuthController {
	return &AuthController{service: service}
}

func (*AuthController) Authenticate(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Authenticate")
}

func (*AuthController) RefreshToken(ctx *gin.Context) {
	ctx.String(http.StatusOK, "RefreshToken")
}

func (*AuthController) LogOut(ctx *gin.Context) {
	ctx.String(http.StatusOK, "LogOut")
}
