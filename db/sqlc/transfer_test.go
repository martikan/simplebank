package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/martikan/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T) Transfer {

	args := CreateTransferParams{
		FromAccountID: 1,
		ToAccountID:   2,
		Amount:        util.RandomUtils.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, args.FromAccountID, transfer.FromAccountID)
	require.Equal(t, args.ToAccountID, transfer.ToAccountID)
	require.Equal(t, args.Amount, transfer.Amount)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestListTransfers(t *testing.T) {

	for i := 0; i < 10; i++ {
		createRandomTransfer(t)
	}

	args := ListTransfersParams{
		Limit:  5,
		Offset: 5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), args)

	require.NoError(t, err)
	require.Len(t, transfers, 5)
	for _, acc := range transfers {
		require.NotEmpty(t, acc)
	}
}

func TestGetTransfer(t *testing.T) {

	transfer1 := createRandomTransfer(t)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transfer2)
	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestUpdateTransfer(t *testing.T) {

	transfer1 := createRandomTransfer(t)

	args := UpdateTransferParams{
		ID:     transfer1.ID,
		Amount: util.RandomUtils.RandomMoney(),
	}

	transfer2, err := testQueries.UpdateTransfer(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, transfer2)
	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, args.Amount, transfer2.Amount)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

func TestDeleteTransfer(t *testing.T) {

	transfer1 := createRandomTransfer(t)

	delErr := testQueries.DeleteTransfer(context.Background(), transfer1.ID)
	transfer2, getErr := testQueries.GetTransfer(context.Background(), transfer1.ID)

	require.NoError(t, delErr)
	require.EqualError(t, getErr, sql.ErrNoRows.Error())
	require.Empty(t, transfer2)
}
