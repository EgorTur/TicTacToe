package service

import (
	"context"
	"testing"
	"tic-tac-toe/internal/domain/entity"
	"tic-tac-toe/internal/domain/repository"
	"tic-tac-toe/internal/domain/repository/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestRegister_Success(t *testing.T) {
    mockRepo := new(mocks.UserRepositoryMock)
    svc := NewService(mockRepo)

    ctx := context.Background()
    login := "newuser"
    password := "123456"

    // Проверка, что пользователь не существует
    mockRepo.On("GetByLogin", ctx, login).Return(nil, repository.ErrUserNotFound)
    // Ожидаем вызов Create
    mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.User")).Return(nil)

    err := svc.Register(ctx, login, password)
    assert.NoError(t, err)

    // Проверим, что пароль был захэширован (косвенно, через аргумент Create)
    mockRepo.AssertCalled(t, "Create", ctx, mock.MatchedBy(func(u *entity.User) bool {
        return u.Login == login && bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) == nil
    }))
    mockRepo.AssertExpectations(t)
}

func TestRegister_AlreadyExists(t *testing.T) {
    mockRepo := new(mocks.UserRepositoryMock)
    svc := NewService(mockRepo)

    ctx := context.Background()
    login := "existing"
    mockRepo.On("GetByLogin", ctx, login).Return(&entity.User{}, nil)

    err := svc.Register(ctx, login, "password")
    assert.EqualError(t, err, "user already exists")
    mockRepo.AssertExpectations(t)
}

func TestAuthenticate_Success(t *testing.T) {
    mockRepo := new(mocks.UserRepositoryMock)
    svc := NewService(mockRepo)

    ctx := context.Background()
    login := "user"
    password := "correct"
    hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    user := &entity.User{ID: uuid.New(), Login: login, PasswordHash: string(hash)}

    mockRepo.On("GetByLogin", ctx, login).Return(user, nil)

    id, err := svc.Authenticate(ctx, login, password)
    assert.NoError(t, err)
    assert.Equal(t, user.ID, id)
    mockRepo.AssertExpectations(t)
}

func TestAuthenticate_WrongPassword(t *testing.T) {
    mockRepo := new(mocks.UserRepositoryMock)
    svc := NewService(mockRepo)

    ctx := context.Background()
    login := "user"
    password := "correct"
    hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    user := &entity.User{ID: uuid.New(), Login: login, PasswordHash: string(hash)}

    mockRepo.On("GetByLogin", ctx, login).Return(user, nil)

    _, err := svc.Authenticate(ctx, login, "wrong")
    assert.EqualError(t, err, "invalid login or password")
    mockRepo.AssertExpectations(t)
}

func TestAuthenticate_UserNotFound(t *testing.T) {
    mockRepo := new(mocks.UserRepositoryMock)
    svc := NewService(mockRepo)

    ctx := context.Background()
    login := "nobody"
    mockRepo.On("GetByLogin", ctx, login).Return(nil, repository.ErrUserNotFound)

    _, err := svc.Authenticate(ctx, login, "whatever")
    assert.EqualError(t, err, "invalid login or password")
    mockRepo.AssertExpectations(t)
}

func TestGetByID_Success(t *testing.T) {
    mockRepo := new(mocks.UserRepositoryMock)
    svc := NewService(mockRepo)

    ctx := context.Background()
    id := uuid.New()
    expected := &entity.User{ID: id, Login: "test"}
    mockRepo.On("GetByID", ctx, id).Return(expected, nil)

    user, err := svc.GetByID(ctx, id)
    assert.NoError(t, err)
    assert.Equal(t, expected, user)
    mockRepo.AssertExpectations(t)
}