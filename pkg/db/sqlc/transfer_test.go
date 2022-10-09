package db

import (
	"context"
	"github.com/samuelowad/bank/pkg/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func createTransferFunc(t *testing.T) Transfer {
	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)
	arg := CreateTransferParams{
		FromAccountID: acc1.ID,
		ToAccountID:   acc2.ID,
		Amount:        util.RandMoney(),
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	return transfer
}

func TestCreateTransfer(t *testing.T) {
	tra := createTransferFunc(t)
	require.NotEmpty(t, tra)
}

func TestGetTransfer(t *testing.T) {
	createTransferFunc(t)
	transfer, err := testQueries.GetTransfer(context.Background(), 1)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

}

func TestListTransfer(t *testing.T) {
	arg := ListTransferParams{
		FromAccountID: 45,
		ToAccountID:   46,
		Offset:        10,
	}
	transfers, err := testQueries.ListTransfer(context.Background(), arg)
	require.NoError(t, err)
	for _, tra := range transfers {
		require.NotEmpty(t, tra)
	}
}
