package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/martikan/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {

	args := CreateAccountParams{
		Owner:    util.RandomUtils.RandomOwner(),
		Balance:  util.RandomUtils.RandomMoney(),
		Currency: util.RandomUtils.RandomCurrency(),
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

func TestListAccounts(t *testing.T) {

	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	args := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), args)

	require.NoError(t, err)
	require.Len(t, accounts, 5)
	for _, acc := range accounts {
		require.NotEmpty(t, acc)
	}
}

func TestGetAccount(t *testing.T) {

	acc1 := createRandomAccount(t)

	acc2, err := testQueries.GetAccount(context.Background(), acc1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, acc2)
	require.Equal(t, acc1.ID, acc2.ID)
	require.Equal(t, acc1.Balance, acc2.Balance)
	require.Equal(t, acc1.Owner, acc2.Owner)
	require.Equal(t, acc1.Currency, acc2.Currency)
	require.WithinDuration(t, acc1.CreatedAt, acc2.CreatedAt, time.Second)
}

func TestGetAccountForUpdate(t *testing.T) {

	acc1 := createRandomAccount(t)

	acc2, err := testQueries.GetAccountForUpdate(context.Background(), acc1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, acc2)
	require.Equal(t, acc1.ID, acc2.ID)
	require.Equal(t, acc1.Balance, acc2.Balance)
	require.Equal(t, acc1.Owner, acc2.Owner)
	require.Equal(t, acc1.Currency, acc2.Currency)
	require.WithinDuration(t, acc1.CreatedAt, acc2.CreatedAt, time.Second)
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestUpdateAccount(t *testing.T) {

	acc1 := createRandomAccount(t)

	args := UpdateAccountParams{
		ID:      acc1.ID,
		Balance: util.RandomUtils.RandomMoney(),
	}

	acc2, err := testQueries.UpdateAccount(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, acc2)
	require.Equal(t, acc1.ID, acc2.ID)
	require.Equal(t, args.Balance, acc2.Balance)
	require.Equal(t, acc1.Owner, acc2.Owner)
	require.Equal(t, acc1.Currency, acc2.Currency)
	require.WithinDuration(t, acc1.CreatedAt, acc2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {

	acc1 := createRandomAccount(t)

	delErr := testQueries.DeleteAccount(context.Background(), acc1.ID)
	acc2, getErr := testQueries.GetAccount(context.Background(), acc1.ID)

	require.NoError(t, delErr)
	require.EqualError(t, getErr, sql.ErrNoRows.Error())
	require.Empty(t, acc2)
}
