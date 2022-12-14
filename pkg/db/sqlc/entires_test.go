package db

import (
	"context"
	"github.com/samuelowad/bank/pkg/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateEntry(t *testing.T) {
	arg := CreateEntryParams{
		AccountID: 1,
		Amount:    util.RandMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
}

func TestGetEntry(t *testing.T) {
	entry, err := testQueries.GetEntry(context.Background(), 1)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
}

func TestListEntry(t *testing.T) {

	arg := ListEntryParams{
		AccountID: 1,
		Limit:     10,
		Offset:    0,
	}
	entries, err := testQueries.ListEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entries)
}
