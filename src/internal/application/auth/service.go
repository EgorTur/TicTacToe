package auth

import (
	"context"
	"errors"
	"fmt"
	"tic-tac-toe/internal/application/service"
)

type authService struct {
	userService service.UserService
	jwtProvider *JwtProvider
}

func NewAuthService(userService service.UserService, jwtProvider *JwtProvider) *authService {
	return &authService{
		userService: userService,
		jwtProvider: jwtProvider,
	}
}

func (s *authService) SignUp(ctx context.Context, req SignUpRequest) error {
	if err := validateCredentials(req.Login, req.Password); err != nil {
		return err
	}
	return s.userService.Register(ctx, req.Login, req.Password)

}

func (s *authService) SignIn(ctx context.Context, req JwtRequest) (JwtResponse, error) {
	if err := validateCredentials(req.Login, req.Password); err != nil {
		return JwtResponse{}, errors.New("invalid login or password")
	}

	userID, err := s.userService.Authenticate(ctx, req.Login, req.Password)

	if err != nil {
		return JwtResponse{}, err
	}

	accessToken, err := s.jwtProvider.GenerateAccessToken(userID)

	if err != nil {
		return JwtResponse{}, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := s.jwtProvider.GenerateRefreshToken(userID)

	if err != nil {
		return JwtResponse{}, fmt.Errorf("generate refresh token: %w", err)
	}

	return JwtResponse{
		Type:         "Bearer",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}

func (s *authService) RefreshAccessToken(ctx context.Context, refreshToken string) (JwtResponse, error) {
	userID, err := s.jwtProvider.ValidateRefreshToken(refreshToken)
	if err != nil {
		return JwtResponse{}, errors.New("invalid refresh tokin")
	}
	_, err = s.userService.GetByID(ctx, userID)
	if err != nil {
		return JwtResponse{}, errors.New("user not found")
	}
	newAccess, err := s.jwtProvider.GenerateAccessToken(userID)
	if err != nil {
		return JwtResponse{}, fmt.Errorf("generate access token: %w", err)
	}
	return JwtResponse{
		Type:         "Bearer",
		AccessToken:  newAccess,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) RefreshRefreshToken(ctx context.Context, refreshToken string) (JwtResponse, error) {
	userID, err := s.jwtProvider.ValidateRefreshToken(refreshToken)
	if err != nil {
		return JwtResponse{}, errors.New("invalid refresh tokin")
	}
	_, err = s.userService.GetByID(ctx, userID)
	if err != nil {
		return JwtResponse{}, errors.New("user not found")
	}
	newAccess, err := s.jwtProvider.GenerateAccessToken(userID)
	if err != nil {
		return JwtResponse{}, fmt.Errorf("generate access token: %w", err)
	}
	newRefresh, err := s.jwtProvider.GenerateRefreshToken(userID)
	if err != nil {
		return JwtResponse{}, fmt.Errorf("generate refresh token: %w", err)
	}
	return JwtResponse{
		Type:         "Bearer",
		AccessToken:  newAccess,
		RefreshToken: newRefresh,
	}, nil
}

func validateCredentials(login, password string) error {
	if len(login) < 3 || len(login) > 15 {
		return errors.New("login must be between 3 and 15 characters")
	}
	if len(password) < 6 || len(password) > 30 {
		return errors.New("password must be between 6 and 30 characters")
	}
	return nil
}
