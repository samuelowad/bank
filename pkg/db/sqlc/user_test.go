package db

import (
	"context"
	"github.com/samuelowad/bank/pkg/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, _ := util.HashedPassword(util.RandString(6))
	arg := CreateUserParams{
		Username:       util.RandOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	if err != nil {
		return createRandomUser(t)

	}
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equalf(t, arg.Username, user.Username, "user")
	require.NotZerof(t, user.CreatedAt, "CreatedAt")
	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user, user2)

}
