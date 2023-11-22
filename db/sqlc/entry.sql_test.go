package db

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"simple_bank/util"
	"testing"
	"time"
)

func createRandomEntry(t * testing.T) Entry {
	account := createRandomAccount(t)

	params := CreateEntryParams{
		AccountID: account.ID,
		Amount: util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, params.AccountID, entry.AccountID)
	require.Equal(t, params.Amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	entry := createRandomEntry(t)

	fetchedEntry, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, fetchedEntry)

	require.Equal(t, entry.AccountID, fetchedEntry.AccountID)
	require.Equal(t, entry.Amount, fetchedEntry.Amount)
	require.Equal(t, entry.ID, fetchedEntry.ID)
	require.WithinDuration(t, entry.CreatedAt, fetchedEntry.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomEntry(t)
	}

	params := ListEntriesParams{
		Limit: 5,
		Offset: 5,
	}

	entries, err := testQueries.ListEntries(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, entries)

	require.Len(t, entries, 5)
	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}

func TestUpdateEntry(t *testing.T) {
	entry := createRandomEntry(t)

	params := UpdateEntryParams{
		ID: entry.ID,
		Amount: util.RandomMoney(),
	}

	updatedEntry, err := testQueries.UpdateEntry(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, updatedEntry)

	require.Equal(t, params.ID, updatedEntry.ID)
	require.Equal(t, params.Amount, updatedEntry.Amount)
	require.Equal(t, entry.AccountID, updatedEntry.AccountID)
	require.WithinDuration(t, entry.CreatedAt, updatedEntry.CreatedAt, time.Second)
}

func TestDeleteEntry(t *testing.T) {
	entry := createRandomEntry(t)

	err := testQueries.DeleteEntry(context.Background(), entry.ID)
	require.NoError(t, err)

	fetchedEntry, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, fetchedEntry)
}
