package auth

import (
	"context"
	"errors"
	"testing"
	"tic-tac-toe/internal/domain/entity"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Мок UserService для AuthService
type mockUserService struct {
    mock.Mock
}

func (m *mockUserService) Register(ctx context.Context, login, password string) error {
    args := m.Called(ctx, login, password)
    return args.Error(0)
}

func (m *mockUserService) Authenticate(ctx context.Context, login, password string) (uuid.UUID, error) {
    args := m.Called(ctx, login, password)
    return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *mockUserService) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entity.User), args.Error(1)
}

func TestSignIn_Success(t *testing.T) {
    userSvc := new(mockUserService)
    jwtProvider := NewJwtProvider("access", "refresh", time.Hour, 24*time.Hour)
    authSvc := NewAuthService(userSvc, jwtProvider)

    ctx := context.Background()
    userID := uuid.New()
    req := JwtRequest{Login: "test", Password: "123456"}

    userSvc.On("Authenticate", ctx, req.Login, req.Password).Return(userID, nil)

    resp, err := authSvc.SignIn(ctx, req)
    assert.NoError(t, err)
    assert.Equal(t, "Bearer", resp.Type)
    assert.NotEmpty(t, resp.AccessToken)
    assert.NotEmpty(t, resp.RefreshToken)

    // Проверяем, что токены содержат правильный userID
    extractedID, err := jwtProvider.ValidateAccessToken(resp.AccessToken)
    assert.NoError(t, err)
    assert.Equal(t, userID, extractedID)

    userSvc.AssertExpectations(t)
}

func TestSignIn_InvalidCredentials(t *testing.T) {
    userSvc := new(mockUserService)
    jwtProvider := NewJwtProvider("access", "refresh", time.Hour, time.Hour)
    authSvc := NewAuthService(userSvc, jwtProvider)

    ctx := context.Background()
    req := JwtRequest{Login: "test", Password: "123456"}  // достаточно длинный
    userSvc.On("Authenticate", ctx, req.Login, req.Password).
        Return(uuid.Nil, errors.New("invalid login or password"))

    _, err := authSvc.SignIn(ctx, req)
    assert.EqualError(t, err, "invalid login or password")
    userSvc.AssertExpectations(t)
}

func TestRefreshAccessToken_Success(t *testing.T) {
    userSvc := new(mockUserService)
    jwtProvider := NewJwtProvider("access", "refresh", time.Hour, 24*time.Hour)
    authSvc := NewAuthService(userSvc, jwtProvider)

    ctx := context.Background()
    userID := uuid.New()
    refreshToken, _ := jwtProvider.GenerateRefreshToken(userID)

    // Мок GetByID для проверки существования пользователя
    userSvc.On("GetByID", ctx, userID).Return(&entity.User{ID: userID}, nil)

    resp, err := authSvc.RefreshAccessToken(ctx, refreshToken)
    assert.NoError(t, err)
    assert.NotEmpty(t, resp.AccessToken)
    // Проверяем, что новый access-токен валиден и содержит тот же userID
    extractedID, err := jwtProvider.ValidateAccessToken(resp.AccessToken)
    assert.NoError(t, err)
    assert.Equal(t, userID, extractedID)
    // Refresh-токен должен остаться прежним (по заданию)
    assert.Equal(t, refreshToken, resp.RefreshToken)

    userSvc.AssertExpectations(t)
}

func TestRefreshAccessToken_UserNotFound(t *testing.T) {
    userSvc := new(mockUserService)
    jwtProvider := NewJwtProvider("access", "refresh", time.Hour, 24*time.Hour)
    authSvc := NewAuthService(userSvc, jwtProvider)

    ctx := context.Background()
    userID := uuid.New()
    refreshToken, _ := jwtProvider.GenerateRefreshToken(userID)

    userSvc.On("GetByID", ctx, userID).Return(nil, errors.New("not found"))

    _, err := authSvc.RefreshAccessToken(ctx, refreshToken)
    assert.EqualError(t, err, "user not found")
    userSvc.AssertExpectations(t)
}

func TestRefreshRefreshToken_Success(t *testing.T) {
    userSvc := new(mockUserService)
    jwtProvider := NewJwtProvider("access", "refresh", time.Hour, 24*time.Hour)
    authSvc := NewAuthService(userSvc, jwtProvider)

    ctx := context.Background()
    userID := uuid.New()
    oldRefresh, _ := jwtProvider.GenerateRefreshToken(userID)

    userSvc.On("GetByID", ctx, userID).Return(&entity.User{ID: userID}, nil)

    // Ждём, чтобы гарантировать другую временную метку
    time.Sleep(1 * time.Second)

    resp, err := authSvc.RefreshRefreshToken(ctx, oldRefresh)
    assert.NoError(t, err)
    assert.NotEmpty(t, resp.AccessToken)
    assert.NotEmpty(t, resp.RefreshToken)
    assert.NotEqual(t, oldRefresh, resp.RefreshToken)

    extractedID, err := jwtProvider.ValidateRefreshToken(resp.RefreshToken)
    assert.NoError(t, err)
    assert.Equal(t, userID, extractedID)

    userSvc.AssertExpectations(t)
}