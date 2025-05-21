package tests

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/noueii/gonuxt-starter/internal/token"
	"github.com/noueii/gonuxt-starter/internal/util"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := token.NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomString(10)
	role := util.UserRole
	duration := time.Minute
	issuedAt := time.Now()
	expiresAt := issuedAt.Add(duration)
	userUUID := uuid.New()
	email := util.RandomEmail()

	token, payload, err := maker.CreateToken(userUUID, email, username, role, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.Equal(t, userUUID, payload.UserID)
	require.Equal(t, email, payload.Email)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiresAt, payload.ExpiresAt, time.Second)
}

func TestExpiredJWTToken(t *testing.T) {
	maker, err := token.NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)
	role := util.UserRole
	userUUID := uuid.New()

	token, payload, err := maker.CreateToken(userUUID, util.RandomEmail(), util.RandomString(10), role, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.ErrorIs(t, err, jwt.ErrTokenExpired)
	require.Nil(t, payload)
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	userUUID := uuid.New()
	payload, err := token.NewPayload(userUUID, util.RandomEmail(), util.RandomString(10), util.UserRole, time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	tk, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := token.NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(tk)
	require.Error(t, err)
	require.ErrorIs(t, err, jwt.ErrTokenSignatureInvalid)
	require.Nil(t, payload)
}
