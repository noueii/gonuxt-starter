package util

import (
	"testing"

	"github.com/noueii/gonuxt-starter/internal/util"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := util.RandomString(10)

	hashed, err := util.HashPassword(password)

	require.NoError(t, err)
	require.NotEmpty(t, hashed)

	err = util.CheckPassword(password, hashed)
	require.NoError(t, err)
}

func TestWrongPassword(t *testing.T) {
	password := util.RandomString(10)
	wrong := util.RandomString(10)

	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	hashedWrongPassword, err := util.HashPassword(wrong)
	require.NoError(t, err)
	require.NotEmpty(t, hashedWrongPassword)

	err = util.CheckPassword(password, wrong)
	require.EqualError(t, err, bcrypt.ErrHashTooShort.Error())

	err = util.CheckPassword(password, hashedWrongPassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}
