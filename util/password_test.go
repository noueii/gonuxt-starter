package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(10)

	hashed, err := HashPassword(password)

	require.NoError(t, err)
	require.NotEmpty(t, hashed)

	err = CheckPassword(password, hashed)
	require.NoError(t, err)
}

func TestWrongPassword(t *testing.T) {
	password := RandomString(10)
	wrong := RandomString(10)

	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	hashedWrongPassword, err := HashPassword(wrong)
	require.NoError(t, err)
	require.NotEmpty(t, hashedWrongPassword)

	err = CheckPassword(password, wrong)
	require.EqualError(t, err, bcrypt.ErrHashTooShort.Error())

	err = CheckPassword(password, hashedWrongPassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}
