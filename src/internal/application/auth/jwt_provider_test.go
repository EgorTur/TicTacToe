package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGenerateAccessToken_Valid(t *testing.T) {
	provider := NewJwtProvider("accessSecret", "refresh_secret", 15*time.Minute, 24*time.Hour)
	userID := uuid.New()
	token, err := provider.GenerateAccessToken(userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGenerateRefreshToken_Valid(t *testing.T) {
    provider := NewJwtProvider("access", "refresh", 15*time.Minute, 24*time.Hour)
    userID := uuid.New()
    token, err := provider.GenerateRefreshToken(userID)
    assert.NoError(t, err)
    assert.NotEmpty(t, token)
}

func TestValidateAccessToken_InvalidSignature(t *testing.T) {
    provider := NewJwtProvider("access_secret", "refresh_secret", 15*time.Minute, 24*time.Hour)
    userID := uuid.New()
    token, _ := provider.GenerateAccessToken(userID)
    // Используем другой провайдер с неправильным секретом
    anotherProvider := NewJwtProvider("wrong", "refresh_secret", 15*time.Minute, 24*time.Hour)
    _, err := anotherProvider.ValidateAccessToken(token)
    assert.Error(t, err)
}

func TestValidateAccessToken_Expired(t *testing.T) {
    provider := NewJwtProvider("access_secret", "refresh_secret", -1*time.Minute, 24*time.Hour) // сразу протухает
    userID := uuid.New()
    token, _ := provider.GenerateAccessToken(userID)
    _, err := provider.ValidateAccessToken(token)
    assert.Error(t, err)
}

func TestValidateRefreshToken_WrongSecret(t *testing.T) {
    provider := NewJwtProvider("access", "refresh", 15*time.Minute, 24*time.Hour)
    userID := uuid.New()
    token, _ := provider.GenerateRefreshToken(userID)

    // Пробуем провалидировать с access‑секретом
    providerWithAccessSecret := NewJwtProvider("access", "access", 15*time.Minute, 24*time.Hour)
    _, err := providerWithAccessSecret.ValidateRefreshToken(token)
    assert.Error(t, err)
}

func TestValidateRefreshToken_Expired(t *testing.T) {
    provider := NewJwtProvider("access", "refresh", 15*time.Minute, -1*time.Hour) // refresh TTL отрицательный
    userID := uuid.New()
    token, _ := provider.GenerateRefreshToken(userID)
    _, err := provider.ValidateRefreshToken(token)
    assert.Error(t, err)
}

func TestGetUserIDFromToken_Success(t *testing.T) {
    provider := NewJwtProvider("access", "refresh", 15*time.Minute, 24*time.Hour)
    userID := uuid.New()
    token, _ := provider.GenerateAccessToken(userID)
    extractedID, err := provider.GetUserIDFromToken(token)
    assert.NoError(t, err)
    assert.Equal(t, userID, extractedID)
}