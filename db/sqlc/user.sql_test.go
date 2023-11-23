package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"simple_bank/util"
	"testing"
	"time"
)

func createRandomUser(t *testing.T) User {
	params := CreateUserParams{
		Username: util.RandomOwner(),
		HashedPassword: "secret",
		FullName: util.RandomOwner(),
		Email: util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, params.Username, user.Username)
	require.Equal(t, params.HashedPassword, user.HashedPassword)
	require.Equal(t, params.FullName, user.FullName)
	require.Equal(t, params.Email, user.Email)

	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)
	require.True(t, user.PasswordLastChangedAt.IsZero())

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user := createRandomUser(t)

	fetchedUser, err := testQueries.GetAUser(context.Background(), user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, fetchedUser)

	require.Equal(t, user.ID, fetchedUser.ID)
	require.Equal(t, user.Username, fetchedUser.Username)
	require.Equal(t, user.HashedPassword, fetchedUser.HashedPassword)
	require.Equal(t, user.FullName, fetchedUser.FullName)
	require.Equal(t, user.Email, fetchedUser.Email)

	require.WithinDuration(t, user.CreatedAt, fetchedUser.CreatedAt, time.Second)
	require.WithinDuration(t, user.PasswordLastChangedAt, fetchedUser.PasswordLastChangedAt, time.Second)
}
