package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/ddcoelho/simple_bank/db/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) (Account, Entry) {
	account := createRandomAccount(t)

	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	updateArg := UpdateAccountParams{
		ID:      account.ID,
		Balance: account.Balance + entry.Amount,
	}

	updatedAccount, err := testQueries.UpdateAccount(context.Background(), updateArg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, updatedAccount.Balance, account.Balance+entry.Amount)

	return updatedAccount, entry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	_, newEntry := createRandomEntry(t)

	requestedEntry, err := testQueries.GetEntry(context.Background(), newEntry.ID)

	require.NoError(t, err)
	require.NotEmpty(t, requestedEntry)

	require.Equal(t, newEntry.ID, requestedEntry.ID)
	require.Equal(t, newEntry.Amount, requestedEntry.Amount)
	require.Equal(t, newEntry.AccountID, requestedEntry.AccountID)
	require.WithinDuration(t, newEntry.CreatedAt, requestedEntry.CreatedAt, time.Second)
}

func TestDeleteEntry(t *testing.T) {
	newAccount, newEntry := createRandomEntry(t)

	err := testQueries.DeleteEntry(context.Background(), newEntry.ID)

	require.NoError(t, err)

	deletedEntry, err := testQueries.GetEntry(context.Background(), newEntry.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, deletedEntry)

	arg := UpdateAccountParams{
		ID:      newAccount.ID,
		Balance: newAccount.Balance - newEntry.Amount,
	}

	updatedAccount, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount)

	require.Equal(t, updatedAccount.Balance, newAccount.Balance-newEntry.Amount)
}

func TestListEntries(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomEntry(t)
	}

	arg := ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entries := range entries {
		require.NotEmpty(t, entries)
	}
}
