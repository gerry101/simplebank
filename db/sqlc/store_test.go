package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	// run n concurrent transfer transactions
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccount.ID,
				ToAccountID: toAccount.ID,
				Amount: amount,
			})

			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		 err := <- errs
		 require.NoError(t, err)

		 result := <- results
		 require.NotEmpty(t, result)

		 // check transfer
		 transfer := result.Transfer
		 require.NotEmpty(t, transfer)
		 require.Equal(t, fromAccount.ID, transfer.FromAccountID)
		 require.Equal(t, toAccount.ID, transfer.ToAccountID)
		 require.Equal(t, amount, transfer.Amount)
		 require.NotZero(t, transfer.ID)
		 require.NotZero(t, transfer.CreatedAt)

		 fetchedTransfer, err := store.GetTransfer(context.Background(), transfer.ID)
		 require.NoError(t, err)
		 require.NotEmpty(t, fetchedTransfer)

		 // check from entry
		 fromEntry := result.FromEntry
		 require.NotEmpty(t, fromEntry)
		 require.Equal(t, fromAccount.ID, fromEntry.AccountID)
		 require.Equal(t, -1 * amount, fromEntry.Amount)
		 require.NotZero(t, fromEntry.ID)
		 require.NotZero(t, fromEntry.CreatedAt)

		 fetchedFromEntry, err := store.GetEntry(context.Background(), fromEntry.ID)
		 require.NoError(t, err)
		 require.NotEmpty(t, fetchedFromEntry)

		 // check to entry
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toAccount.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		fetchedToEntry, err := store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)
		require.NotEmpty(t, fetchedToEntry)

		// check from account
		resultFromAccount := result.FromAccount
		require.NotEmpty(t, resultFromAccount)
		require.Equal(t, fromAccount.ID, resultFromAccount.ID)

		// check to account
		resultToAccount := result.ToAccount
		require.NotEmpty(t, resultToAccount)
		require.Equal(t, toAccount.ID, resultToAccount.ID)

		// check accounts' balance
		fromAccountBalanceDelta := fromAccount.Balance - resultFromAccount.Balance
		toAccountBalanceDelta :=  resultToAccount.Balance - toAccount.Balance
		require.Equal(t, fromAccountBalanceDelta, toAccountBalanceDelta)
		require.True(t, fromAccountBalanceDelta > 0)
		require.True(t, fromAccountBalanceDelta % amount == 0)

		k := int(fromAccountBalanceDelta / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check final updated balances
	updatedFromAccount, err := store.GetAccount(context.Background(), fromAccount.ID)
	require.NoError(t, err)
	updatedToAccount, err := store.GetAccount(context.Background(), toAccount.ID)
	require.NoError(t, err)

	require.Equal(t, fromAccount.Balance - (int64(n) * amount), updatedFromAccount.Balance)
	require.Equal(t, toAccount.Balance + (int64(n) * amount), updatedToAccount.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	// run n concurrent transfer transactions
	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := fromAccount.ID
		toAccountID := toAccount.ID

		if i % 2 == 0 {
			fromAccountID = toAccount.ID
			toAccountID = fromAccount.ID
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID: toAccountID,
				Amount: amount,
			})

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <- errs
		require.NoError(t, err)
	}

	// check final updated balances
	updatedFromAccount, err := store.GetAccount(context.Background(), fromAccount.ID)
	require.NoError(t, err)
	updatedToAccount, err := store.GetAccount(context.Background(), toAccount.ID)
	require.NoError(t, err)

	require.Equal(t, fromAccount.Balance, updatedFromAccount.Balance)
	require.Equal(t, toAccount.Balance, updatedToAccount.Balance)
}
