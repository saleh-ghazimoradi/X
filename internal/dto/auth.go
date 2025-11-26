package dto

import (
	"fmt"
	"github.com/saleh-ghazimoradi/X/internal/customErr"
	"github.com/saleh-ghazimoradi/X/internal/domain"
	"regexp"
	"strings"
)

const (
	UsernameMinLength = 2
	PasswordMinLength = 6
)

var emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type AuthenticationInput struct {
	Email           string `json:"email"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type AuthenticationResponse struct {
	AccessToken string       `json:"access_token"`
	User        *domain.User `json:"user"`
}

func (a *AuthenticationInput) Sanitize() {
	a.Email = strings.TrimSpace(a.Email)
	a.Email = strings.ToLower(a.Email)
	a.Username = strings.TrimSpace(a.Username)
}

func (a *AuthenticationInput) Validate() error {
	if len(a.Username) < UsernameMinLength {
		return fmt.Errorf("%w: username not long enough, (%d) character as least", customErr.ErrValidation, UsernameMinLength)
	}

	if !emailRegexp.MatchString(a.Email) {
		return fmt.Errorf("%w: invalid email address", customErr.ErrValidation)
	}

	if len(a.Password) < PasswordMinLength {
		return fmt.Errorf("%w: password not long enough, (%d) character as least", customErr.ErrValidation, PasswordMinLength)
	}

	if a.Password != a.ConfirmPassword {
		return fmt.Errorf("%w: confirm password must match the password", customErr.ErrValidation)
	}
	return nil
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (l *Login) Sanitize() {
	l.Email = strings.TrimSpace(l.Email)
	l.Email = strings.ToLower(l.Email)
}

func (l *Login) Validate() error {
	if !emailRegexp.MatchString(l.Email) {
		return fmt.Errorf("%w: email not valid", customErr.ErrValidation)
	}
	if len(l.Password) < 1 {
		return fmt.Errorf("%w: password required", customErr.ErrValidation)
	}
	return nil
}
