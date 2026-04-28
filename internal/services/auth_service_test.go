package services_test

import (
	"errors"
	"testing"
	"time"

	"github.com/m0xyu/learning-go-shop/internal/config"
	"github.com/m0xyu/learning-go-shop/internal/dto"
	"github.com/m0xyu/learning-go-shop/internal/mocks"
	"github.com/m0xyu/learning-go-shop/internal/models"
	"github.com/m0xyu/learning-go-shop/internal/notifications"
	"github.com/m0xyu/learning-go-shop/internal/services"
	"github.com/m0xyu/learning-go-shop/internal/utils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAuthService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	mockCartRepo := mocks.NewMockCartRepositoryInterface(ctrl)
	mockEventPublisher := mocks.NewMockPublisher(ctrl)

	dummyJWTConfig := &config.JWTConfig{
		Secret:              "secret",
		ExpiresIn:           time.Hour,
		RefreshTokenExpires: time.Hour * 24,
	}

	dummyConfig := &config.Config{
		JWT: *dummyJWTConfig,
	}

	authService := services.NewAuthService(
		mockUserRepo,
		mockCartRepo,
		dummyConfig,
		mockEventPublisher)

	req := &dto.RegisterRequest{
		Email:     "test@test.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		Phone:     "1234567890",
	}

	mockUserRepo.EXPECT().
		GetByEmail(req.Email).
		Return(nil, errors.New("record not found")).
		Times(1)

	mockUserRepo.EXPECT().
		Create(gomock.Any()). // 引数は models.User が渡ってくる
		Return(nil).
		Times(1)

	mockCartRepo.EXPECT().
		Create(gomock.Any()).
		Return(nil).
		Times(1)

	mockUserRepo.EXPECT().
		CreateRefreshToken(gomock.Any()). // 引数は models.RefreshToken が渡ってくる
		Return(nil).
		Times(1)

	mockEventPublisher.EXPECT().
		Publish(
			notifications.UserLoggedIn,
			gomock.Any(),
			gomock.Any(),
		).
		Return(nil).
		Times(1)

	resp, err := authService.Register(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.Email, resp.User.Email)
	assert.Equal(t, req.FirstName, resp.User.FirstName)
	assert.Equal(t, req.LastName, resp.User.LastName)
	assert.Equal(t, req.Phone, resp.User.Phone)
}

func TestAuthService_Register_EmailAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	mockCartRepo := mocks.NewMockCartRepositoryInterface(ctrl)
	mockEventPublisher := mocks.NewMockPublisher(ctrl)

	dummyJWTConfig := &config.JWTConfig{
		Secret:              "secret",
		ExpiresIn:           time.Hour,
		RefreshTokenExpires: time.Hour * 24,
	}

	dummyConfig := &config.Config{
		JWT: *dummyJWTConfig,
	}

	authService := services.NewAuthService(
		mockUserRepo,
		mockCartRepo,
		dummyConfig,
		mockEventPublisher)

	req := &dto.RegisterRequest{
		Email:     "test@test.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		Phone:     "1234567890",
	}

	mockUserRepo.EXPECT().
		GetByEmail(req.Email).
		Return(&models.User{Email: req.Email}, nil).
		Times(1)

	resp, err := authService.Register(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "email already exists", err.Error())
}

func TestAuthService_Register_UserCreationFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	mockCartRepo := mocks.NewMockCartRepositoryInterface(ctrl)
	mockEventPublisher := mocks.NewMockPublisher(ctrl)

	dummyConfig := &config.Config{}

	authService := services.NewAuthService(
		mockUserRepo,
		mockCartRepo,
		dummyConfig,
		mockEventPublisher)

	req := &dto.RegisterRequest{
		Email:     "test@test.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		Phone:     "1234567890",
	}

	mockUserRepo.EXPECT().
		GetByEmail(req.Email).
		Return(nil, errors.New("record not found")).
		Times(1)

	dbError := errors.New("database connection lost")
	mockUserRepo.EXPECT().
		Create(gomock.Any()).
		Return(dbError).
		Times(1)

	resp, err := authService.Register(req)

	assert.Error(t, err)
	assert.EqualError(t, err, "database connection lost")
	assert.Nil(t, resp)
}

func TestAuthService_Register_CartCreationFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	mockCartRepo := mocks.NewMockCartRepositoryInterface(ctrl)
	mockEventPublisher := mocks.NewMockPublisher(ctrl)

	dummyConfig := &config.Config{}

	authService := services.NewAuthService(
		mockUserRepo,
		mockCartRepo,
		dummyConfig,
		mockEventPublisher)

	req := &dto.RegisterRequest{
		Email:     "test@test.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		Phone:     "1234567890",
	}

	mockUserRepo.EXPECT().
		GetByEmail(req.Email).
		Return(nil, errors.New("record not found")).
		Times(1)

	mockUserRepo.EXPECT().
		Create(gomock.Any()).
		Return(nil).
		Times(1)

	cartDbError := errors.New("failed to connect to cart db")
	mockCartRepo.EXPECT().
		Create(gomock.Any()).
		Return(cartDbError).
		Times(1)

	mockUserRepo.EXPECT().
		CreateRefreshToken(gomock.Any()).
		Return(nil).
		Times(1)

	mockEventPublisher.EXPECT().
		Publish(notifications.UserLoggedIn, gomock.Any(), gomock.Any()).
		Return(nil).
		Times(1)

	resp, err := authService.Register(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.Email, resp.User.Email)
}

func TestAuthService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	mockCartRepo := mocks.NewMockCartRepositoryInterface(ctrl)
	mockEventPublisher := mocks.NewMockPublisher(ctrl)

	dummyJWTConfig := &config.JWTConfig{
		Secret:              "secret",
		ExpiresIn:           time.Hour,
		RefreshTokenExpires: time.Hour * 24,
	}

	dummyConfig := &config.Config{
		JWT: *dummyJWTConfig,
	}

	authService := services.NewAuthService(
		mockUserRepo,
		mockCartRepo,
		dummyConfig,
		mockEventPublisher)

	req := &dto.LoginRequest{
		Email:    "test@test.com",
		Password: "password123",
	}

	hashedPassword, _ := utils.HashPassword(req.Password)

	mockUserRepo.EXPECT().
		GetByEmailAndActive(req.Email, true).
		Return(&models.User{
			ID:       1,
			Email:    req.Email,
			Password: hashedPassword,
			Role:     models.UserRoleCustomer,
		}, nil).
		Times(1)

	mockUserRepo.EXPECT().
		CreateRefreshToken(gomock.Any()).
		Return(nil).
		Times(1)

	mockEventPublisher.EXPECT().
		Publish(
			notifications.UserLoggedIn,
			gomock.Any(),
			gomock.Any(),
		).
		Return(nil).
		Times(1)

	resp, err := authService.Login(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.Email, resp.User.Email)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	mockCartRepo := mocks.NewMockCartRepositoryInterface(ctrl)
	mockEventPublisher := mocks.NewMockPublisher(ctrl)

	dummyConfig := &config.Config{}

	authService := services.NewAuthService(
		mockUserRepo,
		mockCartRepo,
		dummyConfig,
		mockEventPublisher)

	req := &dto.LoginRequest{
		Email:    "test@test.com",
		Password: "password123",
	}

	mockUserRepo.EXPECT().
		GetByEmailAndActive(req.Email, true).
		Return(nil, errors.New("record not found")).
		Times(1)

	resp, err := authService.Login(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "invalid credentials", err.Error())
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	mockCartRepo := mocks.NewMockCartRepositoryInterface(ctrl)
	mockEventPublisher := mocks.NewMockPublisher(ctrl)

	dummyConfig := &config.Config{}

	authService := services.NewAuthService(
		mockUserRepo,
		mockCartRepo,
		dummyConfig,
		mockEventPublisher)

	req := &dto.LoginRequest{
		Email:    "test@test.com",
		Password: "password123",
	}

	hashedPassword, _ := utils.HashPassword("wrong_password")

	mockUserRepo.EXPECT().
		GetByEmailAndActive(req.Email, true).
		Return(&models.User{
			ID:       1,
			Email:    req.Email,
			Password: hashedPassword,
			Role:     models.UserRoleCustomer,
		}, nil).
		Times(1)

	resp, err := authService.Login(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "invalid credentials", err.Error())
}

func TestAuthService_RefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	mockCartRepo := mocks.NewMockCartRepositoryInterface(ctrl)
	mockEventPublisher := mocks.NewMockPublisher(ctrl)

	dummyJWTConfig := &config.JWTConfig{
		Secret:              "secret",
		ExpiresIn:           time.Hour,
		RefreshTokenExpires: time.Hour * 24,
	}

	dummyConfig := &config.Config{
		JWT: *dummyJWTConfig,
	}

	authService := services.NewAuthService(
		mockUserRepo,
		mockCartRepo,
		dummyConfig,
		mockEventPublisher)

	userID := uint(1)
	userEmail := "test@test.com"
	_, realRefreshToken, err := utils.GenerateTokenPair(&dummyConfig.JWT, userID, userEmail, "CUSTOMER")
	assert.NoError(t, err)

	req := &dto.RefreshTokenRequest{
		RefreshToken: realRefreshToken,
	}

	mockUserRepo.EXPECT().
		GetValidRefreshToken(req.RefreshToken).
		Return(&models.RefreshToken{UserID: userID}, nil).
		Times(1)

	mockUserRepo.EXPECT().
		GetByID(userID).
		Return(&models.User{ID: userID, Email: userEmail, Role: models.UserRoleCustomer}, nil).
		Times(1)

	mockUserRepo.EXPECT().
		DeleteRefreshTokenByID(gomock.Any()). // 引数は uint 型のユーザーIDが渡ってくる
		Return(nil).
		Times(1)

	mockUserRepo.EXPECT().
		CreateRefreshToken(gomock.Any()). // 引数は models.RefreshToken が渡ってくる
		Return(nil).
		Times(1)

	mockEventPublisher.EXPECT().
		Publish(
			notifications.UserLoggedIn,
			gomock.Any(),
			gomock.Any(),
		).
		Return(nil).
		Times(1)

	resp, err := authService.RefreshToken(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
}

func TestAuthService_RefreshToken_InvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	mockCartRepo := mocks.NewMockCartRepositoryInterface(ctrl)
	mockEventPublisher := mocks.NewMockPublisher(ctrl)

	dummyConfig := &config.Config{}

	authService := services.NewAuthService(
		mockUserRepo,
		mockCartRepo,
		dummyConfig,
		mockEventPublisher)

	req := &dto.RefreshTokenRequest{
		RefreshToken: "invalid_token",
	}

	resp, err := authService.RefreshToken(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "invalid refresh token", err.Error())
}

func TestAuthService_RefreshToken_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	mockCartRepo := mocks.NewMockCartRepositoryInterface(ctrl)
	mockEventPublisher := mocks.NewMockPublisher(ctrl)

	dummyJWTConfig := &config.JWTConfig{
		Secret:              "secret",
		ExpiresIn:           time.Hour,
		RefreshTokenExpires: time.Hour * 24,
	}

	dummyConfig := &config.Config{
		JWT: *dummyJWTConfig,
	}

	authService := services.NewAuthService(
		mockUserRepo,
		mockCartRepo,
		dummyConfig,
		mockEventPublisher)

	userID := uint(1)
	_, realRefreshToken, err := utils.GenerateTokenPair(
		&dummyConfig.JWT,
		userID,
		"test@test.com",
		string(models.UserRoleCustomer))
	assert.NoError(t, err)

	req := &dto.RefreshTokenRequest{
		RefreshToken: realRefreshToken,
	}

	mockUserRepo.EXPECT().
		GetValidRefreshToken(req.RefreshToken).
		Return(&models.RefreshToken{UserID: userID}, nil).
		Times(1)

	mockUserRepo.EXPECT().
		GetByID(userID).
		Return(nil, errors.New("record not found")).
		Times(1)

	resp, err := authService.RefreshToken(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "user not found", err.Error())
}

func TestAuthService_Logout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	mockCartRepo := mocks.NewMockCartRepositoryInterface(ctrl)
	mockEventPublisher := mocks.NewMockPublisher(ctrl)

	dummyConfig := &config.Config{}

	authService := services.NewAuthService(
		mockUserRepo,
		mockCartRepo,
		dummyConfig,
		mockEventPublisher)

	refreshToken := "some_refresh_token"

	mockUserRepo.EXPECT().
		DeleteRefreshToken(refreshToken).
		Return(nil).
		Times(1)

	err := authService.Logout(refreshToken)

	assert.NoError(t, err)
}
