package service

import (
	"golang.org/x/crypto/bcrypt"
	"os"
	"testing"
)

func TestMain(t *testing.M) {
	passwordCost = bcrypt.MinCost
	os.Exit(t.Run())
}
