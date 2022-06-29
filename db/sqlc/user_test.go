package db

import (
	"context"
	"testing"
	"time"

	"github.com/martikan/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {

	hashedPass, err := util.PasswordUtils.HashPassword(util.RandomUtils.RandomString(6))
	require.NoError(t, err)

	args := CreateUserParams{
		Username: util.RandomUtils.RandomOwner(),
		Email:    util.RandomUtils.RandomEmail(),
		Password: hashedPass,
		FullName: util.RandomUtils.RandomOwner(),
	}

	user, err := testQueries.CreateUser(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, args.Username, user.Username)
	require.Equal(t, args.Email, user.Email)
	require.Equal(t, args.Password, user.Password)
	require.Equal(t, args.FullName, user.FullName)
	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)
	require.True(t, user.PasswordChangedAt.IsZero())

	return user
}

func TestGetUser(t *testing.T) {

	usr1 := createRandomUser(t)

	usr2, err := testQueries.GetUser(context.Background(), usr1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, usr2)
	require.Equal(t, usr1.ID, usr2.ID)
	require.Equal(t, usr1.Email, usr2.Email)
	require.Equal(t, usr1.Username, usr2.Username)
	require.Equal(t, usr1.Password, usr2.Password)
	require.Equal(t, usr1.FullName, usr2.FullName)
	require.WithinDuration(t, usr1.CreatedAt, usr2.CreatedAt, time.Second)
	require.WithinDuration(t, usr1.PasswordChangedAt, usr2.PasswordChangedAt, time.Second)
}

func TestGetUserByUsername(t *testing.T) {

	usr1 := createRandomUser(t)

	usr2, err := testQueries.GetUserByUsername(context.Background(), usr1.Username)

	require.NoError(t, err)
	require.NotEmpty(t, usr2)
	require.Equal(t, usr1.ID, usr2.ID)
	require.Equal(t, usr1.Email, usr2.Email)
	require.Equal(t, usr1.Username, usr2.Username)
	require.Equal(t, usr1.Password, usr2.Password)
	require.Equal(t, usr1.FullName, usr2.FullName)
	require.WithinDuration(t, usr1.CreatedAt, usr2.CreatedAt, time.Second)
	require.WithinDuration(t, usr1.PasswordChangedAt, usr2.PasswordChangedAt, time.Second)
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}
