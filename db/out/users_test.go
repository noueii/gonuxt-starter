package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/noueii/gonuxt-starter/util"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func createRandomUser(t *testing.T) User {
	name := util.RandomString(10)
	user, err := testQueries.CreateUser(context.Background(), name)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, name, user.Name)
	require.Equal(t, int32(0), user.Balance)
	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.UpdatedAt)

	return user
}

func TestGetUserById(t *testing.T) {
	createdUser := createRandomUser(t)
	existingUser, err := testQueries.GetUserById(context.Background(), createdUser.ID)
	require.NoError(t, err)
	require.NotEmpty(t, existingUser)

	require.Equal(t, createdUser, existingUser)

}

func TestUpdateUserBalance(t *testing.T) {
	createdUser := createRandomUser(t)
	newBalance := int32(util.RandomInt(5, 1000))

	updatedUser, err := testQueries.UpdateUserBalance(context.Background(), UpdateUserBalanceParams{
		ID:      createdUser.ID,
		Balance: newBalance,
	})

	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)

	require.Equal(t, newBalance, updatedUser.Balance)
}

func TestDeleteUser(t *testing.T) {
	createdUser := createRandomUser(t)

	err := testQueries.DeleteUser(context.Background(), createdUser.ID)
	require.NoError(t, err)

	notExistingUser, err := testQueries.GetUserById(context.Background(), createdUser.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, notExistingUser)
}
