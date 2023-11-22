package db

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"simple_bank/util"
	"testing"
	"time"
)

func createRandomTransfer(t *testing.T) Transfer {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	params := CreateTransferParams{
		FromAccountID: fromAccount.ID,
		ToAccountID: toAccount.ID,
		Amount: util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, params.FromAccountID, transfer.FromAccountID)
	require.Equal(t, params.ToAccountID, transfer.ToAccountID)
	require.Equal(t, params.Amount, transfer.Amount)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	transfer := createRandomTransfer(t)

	fetchedTransfer, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, fetchedTransfer)

	require.Equal(t, transfer.FromAccountID, fetchedTransfer.FromAccountID)
	require.Equal(t, transfer.ToAccountID, fetchedTransfer.ToAccountID)
	require.Equal(t, transfer.Amount, fetchedTransfer.Amount)
	require.Equal(t, transfer.ID, fetchedTransfer.ID)
	require.WithinDuration(t, transfer.CreatedAt, fetchedTransfer.CreatedAt, time.Second)
}

func TestListTransfers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomTransfer(t)
	}

	params := ListTransfersParams{
		Limit: 5,
		Offset: 5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, transfers)

	require.Len(t, transfers, 5)
	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}

func TestUpdateTransfer(t *testing.T) {
	transfer := createRandomTransfer(t)
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	params := UpdateTransferParams{
		ID: transfer.ID,
		FromAccountID: fromAccount.ID,
		ToAccountID: toAccount.ID,
		Amount: util.RandomMoney(),
	}

	updatedTransfer, err := testQueries.UpdateTransfer(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, updatedTransfer)

	require.Equal(t, params.ID, updatedTransfer.ID)
	require.Equal(t, params.FromAccountID, updatedTransfer.FromAccountID)
	require.Equal(t, params.ToAccountID, updatedTransfer.ToAccountID)
	require.Equal(t, params.Amount, updatedTransfer.Amount)
	require.WithinDuration(t, transfer.CreatedAt, updatedTransfer.CreatedAt, time.Second)
}

func TestDeleteTransfer(t *testing.T) {
	transfer := createRandomTransfer(t)

	err := testQueries.DeleteTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)

	fetchedTransfer, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, fetchedTransfer)
}
