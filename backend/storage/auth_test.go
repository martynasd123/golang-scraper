package storage

import (
	"testing"
	"time"
)

func TestCreateUser(t *testing.T) {
	dao := CreateAuthInMemoryDao()
	user := &User{
		Username: "testuser",
		Password: "password",
	}

	err := dao.CreateUser(user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = dao.CreateUser(user)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestGetUser(t *testing.T) {
	dao := CreateAuthInMemoryDao()
	user := &User{
		Username: "testuser",
		Password: "password",
	}

	err := dao.CreateUser(user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	retrievedUser, err := dao.GetUser("testuser")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if retrievedUser.Username != user.Username {
		t.Fatalf("expected username %v, got %v", user.Username, retrievedUser.Username)
	}
	if retrievedUser.Password != user.Password {
		t.Fatalf("expected password %v, got %v", user.Password, retrievedUser.Password)
	}

	_, err = dao.GetUser("nonexistent")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestUpdateUser(t *testing.T) {
	dao := CreateAuthInMemoryDao()
	user := &User{
		Username: "testuser",
		Password: "password",
	}

	err := dao.CreateUser(user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updatedUser := &User{
		Username:               "testuser",
		Password:               "newpassword",
		DeviceIdentifier:       strPtr("newDeviceIdentifier"),
		RefreshToken:           strPtr("newRefreshToken"),
		RefreshTokenValidUntil: timePtr(time.Now().Add(24 * time.Hour)),
	}

	err = dao.UpdateUser(updatedUser)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	retrievedUser, err := dao.GetUser("testuser")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if retrievedUser.Password != updatedUser.Password {
		t.Fatalf("expected password %v, got %v", updatedUser.Password, retrievedUser.Password)
	}
	if *retrievedUser.DeviceIdentifier != *updatedUser.DeviceIdentifier {
		t.Fatalf("expected device identifier %v, got %v", *updatedUser.DeviceIdentifier, *retrievedUser.DeviceIdentifier)
	}
	if *retrievedUser.RefreshToken != *updatedUser.RefreshToken {
		t.Fatalf("expected refresh token %v, got %v", *updatedUser.RefreshToken, *retrievedUser.RefreshToken)
	}
	if !retrievedUser.RefreshTokenValidUntil.Equal(*updatedUser.RefreshTokenValidUntil) {
		t.Fatalf("expected refresh token valid until %v, got %v", *updatedUser.RefreshTokenValidUntil, *retrievedUser.RefreshTokenValidUntil)
	}

	nonExistentUser := &User{
		Username: "nonexistent",
		Password: "password",
	}
	err = dao.UpdateUser(nonExistentUser)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func strPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}
