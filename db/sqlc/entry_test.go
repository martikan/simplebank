package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/martikan/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) Entry {

	args := CreateEntryParams{
		AccountID: 1,
		Amount:    util.RandomUtils.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, args.AccountID, entry.AccountID)
	require.Equal(t, args.Amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestListEntries(t *testing.T) {

	for i := 0; i < 10; i++ {
		createRandomEntry(t)
	}

	args := ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}

	entries, err := testQueries.ListEntries(context.Background(), args)

	require.NoError(t, err)
	require.Len(t, entries, 5)
	for _, acc := range entries {
		require.NotEmpty(t, acc)
	}
}

func TestGetEntry(t *testing.T) {

	entr1 := createRandomEntry(t)

	entr2, err := testQueries.GetEntry(context.Background(), entr1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entr2)
	require.Equal(t, entr1.ID, entr2.ID)
	require.Equal(t, entr1.AccountID, entr2.AccountID)
	require.Equal(t, entr1.Amount, entr2.Amount)
	require.WithinDuration(t, entr1.CreatedAt, entr2.CreatedAt, time.Second)
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestUpdateEntry(t *testing.T) {

	entr1 := createRandomEntry(t)

	args := UpdateEntryParams{
		ID:     entr1.ID,
		Amount: util.RandomUtils.RandomMoney(),
	}

	entr2, err := testQueries.UpdateEntry(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, entr2)
	require.Equal(t, entr1.ID, entr2.ID)
	require.Equal(t, args.Amount, entr2.Amount)
	require.Equal(t, entr1.AccountID, entr2.AccountID)
	require.WithinDuration(t, entr1.CreatedAt, entr2.CreatedAt, time.Second)
}

func TestDeleteEntry(t *testing.T) {

	entr1 := createRandomEntry(t)

	delErr := testQueries.DeleteEntry(context.Background(), entr1.ID)
	entr2, getErr := testQueries.GetEntry(context.Background(), entr1.ID)

	require.NoError(t, delErr)
	require.EqualError(t, getErr, sql.ErrNoRows.Error())
	require.Empty(t, entr2)
}
