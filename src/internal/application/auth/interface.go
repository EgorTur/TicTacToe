package auth

import (
	"context"
)

type AuthService interface {
	SignUp(ctx context.Context, req SignUpRequest) error
	SignIn(ctx context.Context, req JwtRequest) (JwtResponse, error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (JwtResponse, error)
	RefreshRefreshToken(ctx context.Context, refreshToken string) (JwtResponse, error)
}
