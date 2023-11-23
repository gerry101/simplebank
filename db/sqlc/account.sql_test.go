package db

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"simple_bank/util"
	"testing"
	"time"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)

	params := CreateAccountParams{
		Owner: user.Username,
		Balance: util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, params.Owner, account.Owner)
	require.Equal(t, params.Balance, account.Balance)
	require.Equal(t, params.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account := createRandomAccount(t)

	fetchedAccount, err := testQueries.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.NotEmpty(t, fetchedAccount)

	require.Equal(t, account.ID, fetchedAccount.ID)
	require.Equal(t, account.Owner, fetchedAccount.Owner)
	require.Equal(t, account.Balance, fetchedAccount.Balance)
	require.Equal(t, account.Currency, fetchedAccount.Currency)

	require.WithinDuration(t, account.CreatedAt, fetchedAccount.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account := createRandomAccount(t)

	params := UpdateAccountParams{
		ID: account.ID,
		Balance: util.RandomMoney(),
	}

	updatedAccount, err := testQueries.UpdateAccount(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount)

	require.Equal(t, account.ID, updatedAccount.ID)
	require.Equal(t, account.Owner, updatedAccount.Owner)
	require.Equal(t, params.Balance, updatedAccount.Balance)
	require.Equal(t, account.Currency, updatedAccount.Currency)

	require.WithinDuration(t, account.CreatedAt, updatedAccount.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	fetchedAccount, err := testQueries.GetAccount(context.Background(), account.ID)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, fetchedAccount)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	params := ListAccountsParams{
		Limit: 5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
