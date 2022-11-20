package db

import (
	"context"
	"testing"

	"github.com/simplebank/db/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, account1, account2 Account) Transfer {
	args := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, args.Amount, transfer.Amount)
	require.Equal(t, args.FromAccountID, transfer.FromAccountID)
	require.Equal(t, args.ToAccountID, transfer.ToAccountID)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	createRandomTransfer(t, account1, account2)
}

func TestGetTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	transfer := createRandomTransfer(t, account1, account2)

	createdTransfer, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.Equal(t, account1.ID, transfer.FromAccountID)
	require.Equal(t, account2.ID, transfer.ToAccountID)

	require.NoError(t, err)
	require.NotEmpty(t, createdTransfer)

	require.Equal(t, transfer.ID, createdTransfer.ID)
	require.Equal(t, transfer.CreatedAt, createdTransfer.CreatedAt)
}

func TestListTransfer(t *testing.T) {
	const n = 5

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	for i := 0; i < n; i++ {
		createRandomTransfer(t, account1, account2)
	}

	args := ListTransfersParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Limit:         n,
		Offset:        0,
	}
	transfers, err := testQueries.ListTransfers(context.Background(), args)
	require.NoError(t, err)
	require.EqualValues(t, len(transfers), args.Limit)

	for _, entry := range transfers {
		require.NotEmpty(t, entry)
	}
}
