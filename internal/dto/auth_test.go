package dto

import (
	"github.com/saleh-ghazimoradi/X/internal/customErr"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRegisterInput_Sanitize(t *testing.T) {
	input := AuthenticationInput{
		Username:        "  bob  ",
		Email:           "  BOB@gmail.com  ",
		Password:        "password",
		ConfirmPassword: "password",
	}

	want := AuthenticationInput{
		Username:        "bob",
		Email:           "bob@gmail.com",
		Password:        "password",
		ConfirmPassword: "password",
	}

	input.Sanitize()

	require.Equal(t, want, input)
}

func TestRegisterInput_Validate(t *testing.T) {
	testCases := []struct {
		name  string
		input AuthenticationInput
		err   error
	}{
		{
			name: "valid",
			input: AuthenticationInput{
				Username:        "bob",
				Email:           "bob@gmail.com",
				Password:        "password",
				ConfirmPassword: "password",
			},
			err: nil,
		},
		{
			name: "invalid email",
			input: AuthenticationInput{
				Username:        "bob",
				Email:           "bob",
				Password:        "password",
				ConfirmPassword: "password",
			},
			err: customErr.ErrValidation,
		},
		{
			name: "too short username",
			input: AuthenticationInput{
				Username:        "b",
				Email:           "bob@gmail.com",
				Password:        "password",
				ConfirmPassword: "password",
			},
			err: customErr.ErrValidation,
		},
		{
			name: "too short password",
			input: AuthenticationInput{
				Username:        "bob",
				Email:           "bob@gmail.com",
				Password:        "pass",
				ConfirmPassword: "pass",
			},
			err: customErr.ErrValidation,
		},
		{
			name: "confirm password not match",
			input: AuthenticationInput{
				Username:        "bob",
				Email:           "bob@gmail.com",
				Password:        "password",
				ConfirmPassword: "notpassword",
			},
			err: customErr.ErrValidation,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			err := tc.input.Validate()
			if tc.err != nil {
				require.ErrorIs(tt, err, tc.err)
			} else {
				require.NoError(tt, err)
			}
		})
	}
}

func TestLoginInput_Sanitize(t *testing.T) {
	input := Login{
		Email:    "  BOB@gmail.com  ",
		Password: "password",
	}

	want := Login{
		Email:    "bob@gmail.com",
		Password: "password",
	}

	input.Sanitize()

	require.Equal(t, want, input)
}

func TestLoginInput_Validate(t *testing.T) {
	testCases := []struct {
		name  string
		input Login
		err   error
	}{
		{
			name: "valid",
			input: Login{
				Email:    "bob@gmail.com",
				Password: "password",
			},
			err: nil,
		},
		{
			name: "invalid email",
			input: Login{
				Email:    "bob",
				Password: "password",
			},
			err: customErr.ErrValidation,
		},
		{
			name: "empty password",
			input: Login{

				Email:    "bob@gmail.com",
				Password: "",
			},
			err: customErr.ErrValidation,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			err := tc.input.Validate()
			if tc.err != nil {
				require.ErrorIs(tt, err, tc.err)
			} else {
				require.NoError(tt, err)
			}
		})
	}
}
