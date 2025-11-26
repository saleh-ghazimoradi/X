package service

import (
	"context"
	"errors"
	"github.com/saleh-ghazimoradi/X/faker"
	"github.com/saleh-ghazimoradi/X/internal/customErr"
	"github.com/saleh-ghazimoradi/X/internal/domain"
	"github.com/saleh-ghazimoradi/X/internal/dto"
	"github.com/saleh-ghazimoradi/X/internal/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAuthService_Register(t *testing.T) {
	validInput := &dto.AuthenticationInput{
		Username:        "bob",
		Email:           "bob@gmail.com",
		Password:        "password",
		ConfirmPassword: "password",
	}
	t.Run("can register", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		userRepository := &mocks.UserRepositoryMock{}

		userRepository.On("GetByUsername", mock.Anything, mock.Anything).Return(nil, customErr.ErrNotFound)
		userRepository.On("GetByEmail", mock.Anything, mock.Anything).Return(nil, customErr.ErrNotFound)
		userRepository.On("Create", mock.Anything, mock.Anything).Return(&domain.User{
			Id:       "123",
			Username: validInput.Username,
			Email:    validInput.Email,
		}, nil)
		service := NewAuthService(userRepository)
		res, err := service.Register(ctx, validInput)

		require.NoError(t, err)
		require.NotEmpty(t, res.AccessToken)
		require.NotEmpty(t, res.User.Id)
		require.NotEmpty(t, res.User.Email)
		require.NotEmpty(t, res.User.Username)

		userRepository.AssertExpectations(t)
	})

	t.Run("username taken", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		userRepository := &mocks.UserRepositoryMock{}

		userRepository.On("GetByUsername", mock.Anything, mock.Anything).Return(nil, nil)
		service := NewAuthService(userRepository)

		_, err := service.Register(ctx, validInput)
		require.ErrorIs(t, err, customErr.ErrUserNameTaken)
		userRepository.AssertNotCalled(t, "Create")
		userRepository.AssertExpectations(t)
	})

	t.Run("email taken", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		userRepository := &mocks.UserRepositoryMock{}
		userRepository.On("GetByUsername", mock.Anything, mock.Anything).Return(nil, customErr.ErrNotFound)
		userRepository.On("GetByEmail", mock.Anything, mock.Anything).Return(nil, nil)
		service := NewAuthService(userRepository)
		_, err := service.Register(ctx, validInput)
		require.ErrorIs(t, err, customErr.ErrEmailTaken)
		userRepository.AssertNotCalled(t, "Create")
		userRepository.AssertExpectations(t)
	})

	t.Run("create error", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		userRepository := &mocks.UserRepositoryMock{}
		userRepository.On("GetByUsername", mock.Anything, mock.Anything).Return(nil, customErr.ErrNotFound)
		userRepository.On("GetByEmail", mock.Anything, mock.Anything).Return(nil, customErr.ErrNotFound)
		userRepository.On("Create", mock.Anything, mock.Anything).Return(nil, errors.New("something"))

		service := NewAuthService(userRepository)
		_, err := service.Register(ctx, validInput)
		require.Error(t, err)
		userRepository.AssertExpectations(t)
	})

	t.Run("invalid input", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		userRepository := &mocks.UserRepositoryMock{}
		service := NewAuthService(userRepository)
		_, err := service.Register(ctx, &dto.AuthenticationInput{})
		require.ErrorIs(t, err, customErr.ErrValidation)
		userRepository.AssertNotCalled(t, "GetByUsername")
		userRepository.AssertNotCalled(t, "GetByEmail")
		userRepository.AssertNotCalled(t, "Create")
		userRepository.AssertExpectations(t)
	})
}

func TestAuthService_Login(t *testing.T) {
	validInput := &dto.Login{
		Email:    "bob@gmail.com",
		Password: "password",
	}

	t.Run("can login", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		userRepository := &mocks.UserRepositoryMock{}

		userRepository.On("GetByEmail", mock.Anything, mock.Anything).Return(&domain.User{
			Email:    validInput.Email,
			Password: faker.Password,
		}, nil)
		service := NewAuthService(userRepository)
		_, err := service.Login(ctx, validInput)
		require.NoError(t, err)
		userRepository.AssertExpectations(t)
	})

	t.Run("wrong password", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		userRepository := &mocks.UserRepositoryMock{}

		userRepository.On("GetByEmail", mock.Anything, mock.Anything).Return(&domain.User{
			Email:    validInput.Email,
			Password: faker.Password,
		}, nil)
		service := NewAuthService(userRepository)
		validInput.Password = "something"
		_, err := service.Login(ctx, validInput)
		require.ErrorIs(t, err, customErr.ErrBadCredential)
		userRepository.AssertExpectations(t)
	})

	t.Run("email not found", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		userRepository := &mocks.UserRepositoryMock{}

		userRepository.On("GetByEmail", mock.Anything, mock.Anything).Return(nil, customErr.ErrNotFound)
		service := NewAuthService(userRepository)
		_, err := service.Login(ctx, validInput)
		require.ErrorIs(t, err, customErr.ErrBadCredential)
		userRepository.AssertExpectations(t)
	})

	t.Run("get user by email error", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		userRepository := &mocks.UserRepositoryMock{}

		userRepository.On("GetByEmail", mock.Anything, mock.Anything).Return(nil, errors.New("something"))
		service := NewAuthService(userRepository)
		_, err := service.Login(ctx, validInput)
		require.Error(t, err)
		userRepository.AssertExpectations(t)
	})

	t.Run("invalid input", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		userRepository := &mocks.UserRepositoryMock{}
		service := NewAuthService(userRepository)

		_, err := service.Login(ctx, &dto.Login{
			Email:    "bob",
			Password: "",
		})
		require.ErrorIs(t, err, customErr.ErrValidation)
		userRepository.AssertExpectations(t)
	})

}
