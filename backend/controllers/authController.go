package authController

import (
	"github.com/gofiber/fiber/v2"
	authService "github.com/martynasd123/golang-scraper/services"
)

type AuthController struct {
	service *authService.AuthService
}

func NewAuthController(service *authService.AuthService) *AuthController {
	return &AuthController{service: service}
}

func (*AuthController) Authenticate(ctx *fiber.Ctx) error {
	return ctx.SendString("about")
}

func (*AuthController) RefreshToken(ctx *fiber.Ctx) error {
	return ctx.SendString("about")
}

func (*AuthController) LogOut(ctx *fiber.Ctx) error {
	return ctx.SendString("about")
}
