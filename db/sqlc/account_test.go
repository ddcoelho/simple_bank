package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/ddcoelho/simple_bank/db/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	newAccount := createRandomAccount(t)
	requestedAccount, err := testQueries.GetAccount(context.Background(), newAccount.ID)

	require.NoError(t, err)
	require.NotEmpty(t, requestedAccount)

	require.Equal(t, newAccount.ID, requestedAccount.ID)
	require.Equal(t, newAccount.Owner, requestedAccount.Owner)
	require.Equal(t, newAccount.Balance, requestedAccount.Balance)
	require.Equal(t, newAccount.Currency, requestedAccount.Currency)
	require.WithinDuration(t, newAccount.CreatedAt, requestedAccount.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	newAccount := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      newAccount.ID,
		Balance: util.RandomMoney(),
	}

	changedAccount, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, changedAccount)

	require.Equal(t, newAccount.ID, changedAccount.ID)
	require.Equal(t, newAccount.Owner, changedAccount.Owner)
	require.Equal(t, arg.Balance, changedAccount.Balance)
	require.Equal(t, newAccount.Currency, changedAccount.Currency)
	require.WithinDuration(t, newAccount.CreatedAt, changedAccount.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	newAccount := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), newAccount.ID)

	require.NoError(t, err)

	deletedAccount, err := testQueries.GetAccount(context.Background(), newAccount.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, deletedAccount)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
