package tests

import (
	"context"
	"database/sql"
	"testing"

	db "github.com/noueii/gonuxt-starter/internal/db/out"
	"github.com/noueii/gonuxt-starter/internal/util"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func createRandomUser(t *testing.T) db.User {
	password := util.RandomString(10)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	params := db.CreateUserParams{
		Name:  util.RandomString(10),
		Email: util.RandomEmail(),
		HashedPassword: sql.NullString{
			String: hashedPassword,
			Valid:  true,
		},
	}

	user, err := testQueries.CreateUser(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, params.Name, user.Name)
	require.Equal(t, int32(0), user.Balance)
	require.Equal(t, params.Email, user.Email)
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

	updatedUser, err := testQueries.UpdateUserBalance(context.Background(), db.UpdateUserBalanceParams{
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

func TestUpdateUserByName(t *testing.T) {
	createdUser := createRandomUser(t)

	newPassword := util.RandomString(10)
	newHashedPassword, err := util.HashPassword(newPassword)
	newUsername := util.RandomString(8)
	require.NoError(t, err)

	newBalance := int32(util.RandomInt(0, 100))

	updatedUser, err := testQueries.UpdateUserById(context.Background(), db.UpdateUserByIdParams{
		ID: createdUser.ID,
		Name: sql.NullString{
			String: newUsername,
			Valid:  true,
		},
		HashedPassword: sql.NullString{
			String: newHashedPassword,
			Valid:  true,
		},
		Balance: sql.NullInt32{
			Int32: newBalance,
			Valid: true,
		},
	})

	require.NoError(t, err)
	require.Equal(t, updatedUser.HashedPassword.String, newHashedPassword)
	require.Equal(t, updatedUser.Balance, newBalance)
	require.Equal(t, updatedUser.Name, newUsername)
}
