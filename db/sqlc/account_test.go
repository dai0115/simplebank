package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	args := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, args.Owner, account.Owner)
	require.Equal(t, args.Balance, account.Balance)
	require.Equal(t, args.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	//require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestListAccount(t *testing.T) {
	const n = 10
	for i := 0; i < n; i++ {
		createRandomAccount(t)
	}

	args := ListAccountParams{
		Limit:  n,
		Offset: 0,
	}
	accounts, err := testQueries.ListAccount(context.Background(), args)
	require.NoError(t, err)
	require.EqualValues(t, len(accounts), args.Limit)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}

func TestUpdateAccount(t *testing.T) {
	account := createRandomAccount(t)
	args := UpdateAccountParams{
		ID:      account.ID,
		Balance: 10000, // original balance must be between 0〜1000. so this must be different
	}
	updatedAccount, err := testQueries.UpdateAccount(context.Background(), args)
	require.NoError(t, err)
	require.Equal(t, account.ID, updatedAccount.ID)
	require.NotEqual(t, account.Balance, updatedAccount.Balance)
}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)
	id := account.ID
	err := testQueries.DeleteAccount(context.Background(), id)
	require.NoError(t, err)

	// check if deleted
	deletedAccount, err := testQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, deletedAccount)

}
