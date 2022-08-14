package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	fmt.Println("acc1 before balance:", acc1.Balance, acc2.Balance)

	//	run n concurrent transactions
	n := 5
	amount := int64(10)
	errs := make(chan error)
	result := make(chan TransferTxResult)
	for i := 0; i < n; i++ {
		go func() {

			arg := TransferTxParams{
				FromAccountID: acc1.ID,
				ToAccountID:   acc2.ID,
				Amount:        amount,
			}
			res, err := store.TransferTx(context.Background(), arg)
			require.NoError(t, err)
			errs <- err
			result <- res
		}()
	}
	//check result
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		res := <-result
		require.NotEmpty(t, res)
		//	check transfer
		transfer := res.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, acc1.ID, transfer.FromAccountID)
		require.Equal(t, acc2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		//	check account entries
		fromEntry, toEntry := res.FromEntry, res.ToEntry
		require.NotEmpty(t, fromEntry)
		require.NotEmpty(t, toEntry)
		require.Equal(t, acc1.ID, fromEntry.AccountID)
		require.Equal(t, acc2.ID, toEntry.AccountID)

		require.Equal(t, -amount, fromEntry.Amount)
		require.Equal(t, amount, toEntry.Amount)
		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		//	TODO check account balances

		fromAccount, toAccount := res.FromAccount, res.ToAccount
		require.NotEmpty(t, fromAccount)
		require.NotEmpty(t, toAccount)
		require.Equal(t, acc1.ID, fromAccount.ID)
		require.Equal(t, acc2.ID, toAccount.ID)

		diff := acc1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - acc2.Balance
		require.Equal(t, diff, diff2)
		require.True(t, diff > 0)
		require.True(t, diff%amount == 0)

		k := int(diff / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true

	}
	//	check the final balance

	updateAccount1, err := testQueries.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)

	updateAccount2, err := testQueries.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)

	fmt.Printf("acc1 after balance: %d, acc2 after balance: %d\n", updateAccount1.Balance, updateAccount2.Balance)

	require.Equal(t, acc1.Balance-int64(n)*amount, updateAccount1.Balance)
	require.Equal(t, acc2.Balance+int64(n)*amount, updateAccount2.Balance)

}

func TestTransferTxDeadLock(t *testing.T) {
	store := NewStore(testDB)

	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	fmt.Println("acc1 before balance:", acc1.Balance, acc2.Balance)

	//	run n concurrent transactions
	n := 10
	amount := int64(10)
	errs := make(chan error)
	for i := 0; i < n; i++ {
		fromAccID := acc1.ID
		toAccID := acc2.ID

		if i%2 == 1 {
			fromAccID = acc2.ID
			toAccID = acc1.ID
		}

		go func() {

			arg := TransferTxParams{
				FromAccountID: fromAccID,
				ToAccountID:   toAccID,
				Amount:        amount,
			}
			_, err := store.TransferTx(context.Background(), arg)
			require.NoError(t, err)
			errs <- err
		}()
	}
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

	}
	//	check the final balance

	updateAccount1, err := testQueries.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)

	updateAccount2, err := testQueries.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)

	fmt.Printf("acc1 after balance: %d, acc2 after balance: %d\n", updateAccount1.Balance, updateAccount2.Balance)

	require.Equal(t, acc1.Balance, updateAccount1.Balance)
	require.Equal(t, acc2.Balance, updateAccount2.Balance)

}
