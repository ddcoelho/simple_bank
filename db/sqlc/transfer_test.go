package db

import (
	"context"
	"testing"

	"github.com/ddcoelho/simple_bank/db/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, fromAccountId int64, toAccountId int64) Transfer {
	fromAccount, err := testQueries.GetAccount(context.Background(), fromAccountId)
	require.NoError(t, err)
	require.NotEmpty(t, fromAccount)

	toAccount, err := testQueries.GetAccount(context.Background(), toAccountId)
	require.NoError(t, err)
	require.NotEmpty(t, toAccount)

	arg := CreateTransferParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	fromArg := UpdateAccountParams{
		ID:      fromAccount.ID,
		Balance: fromAccount.Balance - transfer.Amount,
	}

	fromUpdatedAccount, err := testQueries.UpdateAccount(context.Background(), fromArg)

	require.NoError(t, err)
	require.NotEmpty(t, fromUpdatedAccount)

	require.Equal(t, fromUpdatedAccount.Balance, fromAccount.Balance-transfer.Amount)

	toArg := UpdateAccountParams{
		ID:      toAccount.ID,
		Balance: toAccount.Balance + transfer.Amount,
	}

	toUpdatedAccount, err := testQueries.UpdateAccount(context.Background(), toArg)
	require.NoError(t, err)
	require.NotEmpty(t, toUpdatedAccount)

	require.Equal(t, toUpdatedAccount.Balance, toAccount.Balance+transfer.Amount)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	createRandomTransfer(t, fromAccount.ID, toAccount.ID)
}

func TestGetTransfer(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	newTransfer := createRandomTransfer(t, fromAccount.ID, toAccount.ID)

	requestedTransfer, err := testQueries.GetTransfer(context.Background(), newTransfer.ID)

	require.NoError(t, err)
	require.NotEmpty(t, requestedTransfer)

	require.Equal(t, newTransfer.ID, requestedTransfer.ID)
	require.Equal(t, newTransfer.Amount, requestedTransfer.Amount)
	require.Equal(t, newTransfer.FromAccountID, requestedTransfer.FromAccountID)
	require.Equal(t, newTransfer.ToAccountID, requestedTransfer.ToAccountID)
}

func TestListTransfers(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomTransfer(t, fromAccount.ID, toAccount.ID)
	}

	arg := ListTransfersParams{
		Limit:  5,
		Offset: 5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfers := range transfers {
		require.NotEmpty(t, transfers)
	}
}

func TestListAccountTransfers(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomTransfer(t, fromAccount.ID, toAccount.ID)
	}

	arg := ListAccountTransfersParams{
		FromAccountID: fromAccount.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListAccountTransfers(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfers := range transfers {
		require.NotEmpty(t, transfers)
	}
}
