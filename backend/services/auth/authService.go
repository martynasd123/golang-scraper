package authService

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	auth "github.com/martynasd123/golang-scraper/storage"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrUserNotExist                     = errors.New("user does not exist")
	ErrRefreshTokenNotExist             = errors.New("refresh token does not exist")
	ErrRefreshTokenMismatch             = errors.New("refresh token does not match")
	ErrRefreshTokenExpired              = errors.New("refresh token is expired")
	ErrAccessTokenInvalid               = errors.New("access token is invalid")
	ErrInvalidDeviceIdentifier          = errors.New("invalid device identifier")
	ErrCouldNotGenerateAccessToken      = errors.New("could not generate access token")
	ErrCouldNotGenerateDeviceIdentifier = errors.New("could not generate device identifier")
	ErrUserPersistenceError             = errors.New("could not update user")
	ErrPasswordIncorrect                = errors.New("password is incorrect")
	ErrCouldNotHashPassword             = errors.New("password hashing failed")
)

var JwtSecretKey = "Some very secret key"

type AuthService struct {
	authStorage auth.AuthDao
}

func CreateAuthService(storage auth.AuthDao) *AuthService {
	return &AuthService{
		authStorage: storage,
	}
}

func (authService *AuthService) Login(username, password, ip, userAgent string) (
	accessToken, refreshToken *string,
	refreshTokenValidUntil *time.Time,
	err error,
) {
	user, err := authService.authStorage.GetUser(username)
	if err != nil {
		return nil, nil, nil, errors.Join(ErrUserNotExist, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, nil, nil, errors.Join(ErrPasswordIncorrect, err)
	}

	refreshToken, refreshTokenValidUntil = generateRefreshToken()

	user.DeviceIdentifier, err = generateDeviceIdentifier(ip, userAgent)
	if err != nil {
		return nil, nil, nil, errors.Join(ErrCouldNotGenerateDeviceIdentifier, err)
	}

	user.RefreshToken = refreshToken
	user.RefreshTokenValidUntil = refreshTokenValidUntil
	err = authService.authStorage.UpdateUser(user)

	if err != nil {
		return nil, nil, nil, errors.Join(ErrUserPersistenceError, err)
	}

	accessToken, err = generateAccessToken(username)
	if err != nil {
		return nil, nil, nil, errors.Join(ErrCouldNotGenerateAccessToken, err)
	}
	return accessToken, refreshToken, refreshTokenValidUntil, nil
}

func generateDeviceIdentifier(ip string, agent string) (*string, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(ip+agent), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	passwordStr := string(password)
	return &passwordStr, err
}

func generateRefreshToken() (*string, *time.Time) {
	refreshToken := uuid.New().String()
	validUntil := time.Now().AddDate(0, 0, 30)
	return &refreshToken, &validUntil
}

func (authService *AuthService) CreateFakeUser(username string, password string) error {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Join(err, ErrCouldNotHashPassword)
	}
	user := &auth.User{
		Username: username,
		Password: string(hashedPass),
	}

	// Store the user in the storage
	err = authService.authStorage.CreateUser(user)
	if err != nil {
		return errors.Join(err, ErrUserPersistenceError)
	}

	return nil
}

func (authService *AuthService) RefreshToken(
	token string,
	username string,
	ip string,
	agent string,
) (accessToken, refreshToken *string, refreshTokenExp *time.Time, err error) {
	user, err := authService.authStorage.GetUser(username)
	if err != nil {
		return nil, nil, nil, errors.Join(ErrUserNotExist, err)
	}

	if user.RefreshToken == nil || user.RefreshTokenValidUntil == nil {
		return nil, nil, nil, ErrRefreshTokenNotExist
	}

	if *user.RefreshToken != token {
		return nil, nil, nil, ErrRefreshTokenMismatch
	}

	if (*user.RefreshTokenValidUntil).Compare(time.Now()) < 0 {
		return nil, nil, nil, ErrRefreshTokenExpired
	}

	err = bcrypt.CompareHashAndPassword([]byte(*user.DeviceIdentifier), []byte(ip+agent))
	if err != nil {
		return nil, nil, nil, errors.Join(ErrInvalidDeviceIdentifier, err)
	}

	accessToken, err = generateAccessToken(username)
	if err != nil {
		return nil, nil, nil, errors.Join(ErrCouldNotGenerateAccessToken, err)
	}

	refreshToken, refreshTokenExp = generateRefreshToken()

	user.RefreshToken = refreshToken
	user.RefreshTokenValidUntil = refreshTokenExp

	err = authService.authStorage.UpdateUser(user)
	if err != nil {
		return nil, nil, nil, errors.Join(ErrUserPersistenceError, err)
	}
	return
}

func (authService *AuthService) LogOut(username string) error {
	user, err := authService.authStorage.GetUser(username)
	if err != nil {
		return errors.Join(ErrUserNotExist, err)
	}

	user.RefreshTokenValidUntil = nil
	user.RefreshToken = nil
	user.DeviceIdentifier = nil
	err = authService.authStorage.UpdateUser(user)

	if err != nil {
		return errors.Join(err, ErrUserPersistenceError)
	}
	return nil
}

func (authService *AuthService) ValidateAccessToken(token string) (username string, err error) {
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrAccessTokenInvalid
		}
		return []byte(JwtSecretKey), nil
	})
	if err != nil {
		return "", errors.Join(ErrAccessTokenInvalid, err)
	}

	if exp, ok := claims["exp"].(float64); ok {
		if time.Unix(int64(exp), 0).Before(time.Now()) {
			return "", ErrAccessTokenInvalid
		}
	} else {
		return "", ErrAccessTokenInvalid
	}

	username, ok := claims["username"].(string)
	if !ok {
		return "", ErrAccessTokenInvalid
	}
	return username, nil
}

func generateAccessToken(username string) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})
	accessToken, err := token.SignedString([]byte(JwtSecretKey))
	if err != nil {
		return nil, err
	}
	return &accessToken, nil
}
