package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JwtProvider struct {
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func NewJwtProvider(accessSecret, refreshSecret string, accessTTL, refreshTTL time.Duration) *JwtProvider {
	return &JwtProvider{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
	}
}

func (p *JwtProvider) GenerateAccessToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(p.accessTTL).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(p.accessSecret)
}

func (p *JwtProvider) GenerateRefreshToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(p.refreshTTL).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(p.refreshSecret)
}

func (p *JwtProvider) ValidateAccessToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return p.accessSecret, nil
	})
	if err != nil || !token.Valid {
		return uuid.Nil, errors.New("invalid access token")

	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, errors.New("invalid claims")
	}
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return uuid.Nil, errors.New("user_id not found in token")
	}
	return uuid.Parse(userIDStr)
}

func (p *JwtProvider) ValidateRefreshToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return p.refreshSecret, nil
	})
	if err != nil || !token.Valid {
		return uuid.Nil, errors.New("invalid refresh token")

	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, errors.New("invalid claims")
	}
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return uuid.Nil, errors.New("user_id not found in token")
	}
	return uuid.Parse(userIDStr)
}

func (p *JwtProvider) GetUserIDFromToken(tokenString string) (uuid.UUID, error) {
	return p.ValidateAccessToken(tokenString)
}
