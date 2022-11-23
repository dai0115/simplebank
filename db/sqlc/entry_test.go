package db

import (
	"context"
	"testing"

	"github.com/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, account Account) Entry {
	args := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, args.Amount, entry.Amount)
	require.Equal(t, account.ID, entry.AccountID)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	createRandomEntry(t, account)
}

func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t)
	entry := createRandomEntry(t, account)
	createdEntry, err := testQueries.GetEntry(context.Background(), entry.ID)

	require.NoError(t, err)
	require.NotEmpty(t, createdEntry)

	require.Equal(t, entry.ID, createdEntry.ID)
	require.Equal(t, entry.AccountID, createdEntry.AccountID)
	require.Equal(t, entry.Amount, createdEntry.Amount)
	//require.WithinDuration(t, entry1.CreatedAt, createdEntry.CreatedAt, time.Second)
}

func TestListEntry(t *testing.T) {
	const n = 5

	account := createRandomAccount(t)
	for i := 0; i < n; i++ {
		createRandomEntry(t, account)
	}

	args := ListEntriesParams{
		AccountID: account.ID,
		Limit:     n,
		Offset:    0,
	}
	entries, err := testQueries.ListEntries(context.Background(), args)
	require.NoError(t, err)
	require.EqualValues(t, len(entries), args.Limit)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
