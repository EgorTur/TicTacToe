package web

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"tic-tac-toe/internal/application/auth"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Мок для AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) SignUp(ctx context.Context, req auth.SignUpRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockAuthService) SignIn(ctx context.Context, req auth.JwtRequest) (auth.JwtResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(auth.JwtResponse), args.Error(1)
}

func (m *MockAuthService) RefreshAccessToken(ctx context.Context, refreshToken string) (auth.JwtResponse, error) {
	args := m.Called(ctx, refreshToken)
	return args.Get(0).(auth.JwtResponse), args.Error(1)
}

func (m *MockAuthService) RefreshRefreshToken(ctx context.Context, refreshToken string) (auth.JwtResponse, error) {
	args := m.Called(ctx, refreshToken)
	return args.Get(0).(auth.JwtResponse), args.Error(1)
}

func TestSignUp_Success(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	reqBody := auth.SignUpRequest{Login: "test", Password: "123456"}
	mockService.On("SignUp", mock.Anything, reqBody).Return(nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/sign-up", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.SignUp(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "registered")
	mockService.AssertExpectations(t)
}

func TestSignUp_InvalidJSON(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/sign-up", bytes.NewReader([]byte("bad json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.SignUp(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSignUp_ServiceError(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	reqBody := auth.SignUpRequest{Login: "test", Password: "123456"}
	mockService.On("SignUp", mock.Anything, reqBody).Return(errors.New("user already exists"))

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/sign-up", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.SignUp(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "user already exists")
	mockService.AssertExpectations(t)
}

func TestSignUp_MethodNotAllowed(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/sign-up", nil)
	w := httptest.NewRecorder()

	handler.SignUp(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestSignIn_Success(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	reqBody := auth.JwtRequest{Login: "test", Password: "123456"}
	expectedResp := auth.JwtResponse{
		Type:         "Bearer",
		AccessToken:  "access_token",
		RefreshToken: "refresh_token",
	}
	mockService.On("SignIn", mock.Anything, reqBody).Return(expectedResp, nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/sign-in", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.SignIn(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp auth.JwtResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, expectedResp, resp)
	mockService.AssertExpectations(t)
}

func TestSignIn_InvalidCredentials(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	reqBody := auth.JwtRequest{Login: "test", Password: "wrong"}
	mockService.On("SignIn", mock.Anything, reqBody).Return(auth.JwtResponse{}, errors.New("invalid login or password"))

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/sign-in", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.SignIn(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid login or password")
	mockService.AssertExpectations(t)
}

func TestSignIn_InvalidJSON(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/sign-in", bytes.NewReader([]byte("bad json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.SignIn(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRefreshAccessToken_Success(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	refreshToken := "valid_refresh_token"
	expectedResp := auth.JwtResponse{
		Type:         "Bearer",
		AccessToken:  "new_access",
		RefreshToken: refreshToken,
	}
	mockService.On("RefreshAccessToken", mock.Anything, refreshToken).Return(expectedResp, nil)

	reqBody := auth.RefreshJwtRequest{RefreshToken: refreshToken}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/refresh-access", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.RefreshAccessToken(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp auth.JwtResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, expectedResp, resp)
	mockService.AssertExpectations(t)
}

func TestRefreshAccessToken_InvalidToken(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	refreshToken := "bad_token"
	mockService.On("RefreshAccessToken", mock.Anything, refreshToken).Return(auth.JwtResponse{}, errors.New("invalid refresh token"))

	reqBody := auth.RefreshJwtRequest{RefreshToken: refreshToken}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/refresh-access", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.RefreshAccessToken(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid refresh token")
	mockService.AssertExpectations(t)
}

func TestRefreshRefreshToken_Success(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	refreshToken := "valid_refresh"
	expectedResp := auth.JwtResponse{
		Type:         "Bearer",
		AccessToken:  "new_access",
		RefreshToken: "new_refresh",
	}
	mockService.On("RefreshRefreshToken", mock.Anything, refreshToken).Return(expectedResp, nil)

	reqBody := auth.RefreshJwtRequest{RefreshToken: refreshToken}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/refresh-refresh", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.RefreshRefreshToken(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp auth.JwtResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, expectedResp, resp)
	mockService.AssertExpectations(t)
}