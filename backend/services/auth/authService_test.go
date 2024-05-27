package authService_test

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/martynasd123/golang-scraper/services/auth"
	"github.com/martynasd123/golang-scraper/storage"
)

func TestAuthService_RefreshToken(t *testing.T) {
	inMemoryStorage := storage.CreateAuthInMemoryDao()
	auth := authService.CreateAuthService(inMemoryStorage)

	err := auth.CreateFakeUser("testuser", "password")
	require.NoError(t, err)

	ip := "192.168.1.1"
	userAgent := "Mozilla/5.0"

	_, refreshToken, _, err := auth.Login("testuser", "password", ip, userAgent)
	require.NoError(t, err)

	_, newRefreshToken, _, err := auth.RefreshToken(*refreshToken, "testuser", ip, userAgent)
	require.NoError(t, err)
	require.NotEqual(t, refreshToken, newRefreshToken)

	// Ip is different - should fail
	_, _, _, err = auth.RefreshToken(*newRefreshToken, "testuser", "192.168.1.2", userAgent)
	require.ErrorIs(t, err, authService.ErrInvalidDeviceIdentifier)
}

func TestAuthService_LogOut(t *testing.T) {
	inMemoryStorage := storage.CreateAuthInMemoryDao()
	auth := authService.CreateAuthService(inMemoryStorage)

	err := auth.CreateFakeUser("testuser", "password")
	require.NoError(t, err)

	_, _, _, err = auth.Login("testuser", "password", "127.0.0.1", "Mozilla/5.0")
	require.NoError(t, err)

	err = auth.LogOut("testuser")
	require.NoError(t, err)

	user, err := inMemoryStorage.GetUser("testuser")
	require.NoError(t, err)
	require.Nil(t, user.RefreshTokenValidUntil)
	require.Nil(t, user.DeviceIdentifier)
	require.Nil(t, user.RefreshToken)

	err = auth.LogOut("nonexistentuser")
	require.Error(t, err)
}

func TestAuthService_ValidateAccessToken(t *testing.T) {
	inMemoryStorage := storage.CreateAuthInMemoryDao()
	auth := authService.CreateAuthService(inMemoryStorage)

	err := auth.CreateFakeUser("testuser", "password")
	require.NoError(t, err)

	ip := "192.168.1.1"
	userAgent := "Mozilla/5.0"

	accessToken, _, _, err := auth.Login("testuser", "password", ip, userAgent)
	require.NoError(t, err)

	username, err := auth.ValidateAccessToken(*accessToken)
	require.NoError(t, err)
	require.Equal(t, "testuser", username)

	_, err = auth.ValidateAccessToken("bad token")
	require.Error(t, err)
}
