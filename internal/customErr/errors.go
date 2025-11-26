package customErr

import "errors"

var (
	ErrValidation    = errors.New("validation error")
	ErrNotFound      = errors.New("not found")
	ErrUserNameTaken = errors.New("user name already taken")
	ErrEmailTaken    = errors.New("email already taken")
	ErrBadCredential = errors.New("email/password wrong combination")
)
