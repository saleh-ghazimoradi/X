package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/saleh-ghazimoradi/X/internal/customErr"
	"github.com/saleh-ghazimoradi/X/internal/domain"
	"github.com/saleh-ghazimoradi/X/internal/dto"
	"github.com/saleh-ghazimoradi/X/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var passwordCost = bcrypt.DefaultCost

type AuthService interface {
	Register(ctx context.Context, input *dto.AuthenticationInput) (*dto.AuthenticationResponse, error)
	Login(ctx context.Context, input *dto.Login) (*dto.AuthenticationResponse, error)
}

type authService struct {
	userRepository repository.UserRepository
}

func (a *authService) Register(ctx context.Context, input *dto.AuthenticationInput) (*dto.AuthenticationResponse, error) {
	input.Sanitize()

	if err := input.Validate(); err != nil {
		return nil, err
	}

	if _, err := a.userRepository.GetByUsername(ctx, input.Username); !errors.Is(err, customErr.ErrNotFound) {
		return nil, customErr.ErrUserNameTaken
	}

	if _, err := a.userRepository.GetByEmail(ctx, input.Email); !errors.Is(err, customErr.ErrNotFound) {
		return nil, customErr.ErrEmailTaken
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), passwordCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %v", err)
	}

	user := &domain.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	if err := a.userRepository.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("error creating user: %v", err)
	}

	return &dto.AuthenticationResponse{
		AccessToken: "a token",
		User:        user,
	}, nil
}

func (a *authService) Login(ctx context.Context, input *dto.Login) (*dto.AuthenticationResponse, error) {
	input.Sanitize()
	if err := input.Validate(); err != nil {
		return nil, err
	}

	user, err := a.userRepository.GetByEmail(ctx, input.Email)
	if err != nil {
		switch {
		case errors.Is(err, customErr.ErrNotFound):
			return nil, customErr.ErrBadCredential
		default:
			return nil, err
		}
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, customErr.ErrBadCredential
	}

	return &dto.AuthenticationResponse{
		AccessToken: "a token",
		User:        user,
	}, nil
}

func NewAuthService(userRepository repository.UserRepository) AuthService {
	return &authService{
		userRepository: userRepository,
	}
}
