package authController

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	. "github.com/martynasd123/golang-scraper/models/request"
	authService "github.com/martynasd123/golang-scraper/services/auth"
	. "github.com/martynasd123/golang-scraper/services/auth/constants"
)

type AuthController struct {
	authService *authService.AuthService
}

func CreateAuthController(service *authService.AuthService) *AuthController {
	return &AuthController{authService: service}
}

func (controller *AuthController) Authenticate(ctx *gin.Context) {
	var authRequest AuthRequest

	if err := ctx.ShouldBindJSON(&authRequest); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	accessToken, refreshToken, refreshTokenValidUntil, err := controller.authService.Login(
		authRequest.Username,
		authRequest.Password,
		ctx.ClientIP(),
		ctx.Request.UserAgent(),
	)
	if err != nil {
		if errors.Is(err, authService.ErrUserNotExist) || errors.Is(err, authService.ErrPasswordIncorrect) {
			ctx.String(http.StatusForbidden, "username or password is incorrect")
		} else {
			log.Println(err)
			ctx.String(http.StatusInternalServerError, "something went wrong")
		}
		return
	}
	setRefreshTokenCookie(ctx, refreshToken, refreshTokenValidUntil)
	setAccessTokenCookie(ctx, accessToken)
	ctx.String(http.StatusOK, "")
}

func setAccessTokenCookie(ctx *gin.Context, accessToken *string) {
	ctx.SetCookie(AccessTokenCookieName,
		*accessToken,
		0,
		"/api/",
		"",
		false,
		true,
	)
}

func setRefreshTokenCookie(ctx *gin.Context, refreshToken *string, refreshTokenValidUntil *time.Time) {
	maxAge := 0
	if refreshTokenValidUntil != nil {
		maxAge = int(refreshTokenValidUntil.Sub(time.Now()).Seconds())
	}
	ctx.SetCookie(RefreshTokenCookieName,
		*refreshToken,
		maxAge,
		"/api/auth/refresh-token",
		"",
		false,
		true,
	)
}

func (controller *AuthController) RefreshToken(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie(RefreshTokenCookieName)
	if err != nil {
		ctx.String(http.StatusBadRequest, "refresh token required")
		return
	}

	var refreshTokenRequest RefreshTokenRequest

	if err := ctx.ShouldBindJSON(&refreshTokenRequest); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	accessToken, newRefreshToken, refreshTokenValidUntil, err := controller.authService.RefreshToken(
		refreshToken,
		refreshTokenRequest.Username,
		ctx.ClientIP(),
		ctx.Request.UserAgent(),
	)
	if err != nil {
		switch {
		case errors.Is(err, authService.ErrRefreshTokenMismatch),
			errors.Is(err, authService.ErrRefreshTokenNotExist),
			errors.Is(err, authService.ErrInvalidDeviceIdentifier),
			errors.Is(err, authService.ErrRefreshTokenExpired),
			errors.Is(err, authService.ErrUserNotExist):
			ctx.String(http.StatusForbidden, "could not verify refresh token")
		default:
			log.Println(fmt.Errorf("unexpected error while refreshing token: %w", err))
			ctx.String(http.StatusInternalServerError, "something went wrong")
		}
		return
	}
	setAccessTokenCookie(ctx, accessToken)
	setRefreshTokenCookie(ctx, newRefreshToken, refreshTokenValidUntil)
}

func (controller *AuthController) LogOut(ctx *gin.Context) {
	var logOutRequest LogOutRequest

	if err := ctx.ShouldBindJSON(&logOutRequest); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	err := controller.authService.LogOut(logOutRequest.Username)
	if err != nil {
		switch {
		case errors.Is(err, authService.ErrUserNotExist):
			ctx.String(http.StatusForbidden, "bad credentials")
		default:
			log.Println(fmt.Errorf("unexpected error while logging out: %w", err))
			ctx.String(http.StatusInternalServerError, "something went wrong")
		}
		return
	}
	setAccessTokenCookie(ctx, new(string))
	setRefreshTokenCookie(ctx, new(string), nil)
}
