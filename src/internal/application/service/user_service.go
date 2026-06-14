package service

import (
	"context"
	"errors"
	"fmt"
	"tic-tac-toe/internal/domain/entity"
	"tic-tac-toe/internal/domain/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo repository.UserRepository
}

func NewService(repo repository.UserRepository) *userService {
	return &userService{repo: repo}
}

func (s *userService) Register(ctx context.Context, login, password string) error {
	_, err := s.repo.GetByLogin(ctx, login)
	if err == nil {
		return errors.New("user already exists")
	}
	if !errors.Is(err, repository.ErrUserNotFound) {
		return fmt.Errorf("check login uniqueness: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	user := &entity.User{
		ID:           uuid.New(),
		Login:        login,
		PasswordHash: string(hash),
	}
	return s.repo.Create(ctx, user)

}

func (s *userService) Authenticate(ctx context.Context, login, password string) (uuid.UUID, error) {
	user, err := s.repo.GetByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return uuid.Nil, errors.New("invalid login or password")
		}
		return uuid.Nil, fmt.Errorf("get user: %w", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return uuid.Nil, errors.New("invalid login or password")
	}
	return user.ID, nil
}

func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	return s.repo.GetByID(ctx, id)
}
