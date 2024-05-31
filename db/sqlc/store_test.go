package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)
	
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	fmt.Println("Saldos iniciales", account1.Balance, account2.Balance)

	n := 5
	amount := int64(10)

	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		result, err := store.TransferTx(context.Background(), TransferTxParams{
			FromAccountID: account1.ID,
			ToAccountID:   account2.ID,
			Amount:        amount,
		})
		require.NoError(t, err)

		require.NotEmpty(t, result)
		
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)
		
		_, err = store.GetTransfer(context.Background(), store.db, transfer.ID)
		require.NoError(t, err)
		
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), store.db, fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), store.db, toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.From
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.To
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		fmt.Println("Saldos en iteracion", fromAccount.Balance, toAccount.Balance)
		fmt.Println("Saldos iniciales", account1.Balance, account2.Balance)
		// checks accounts balances
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true

		fmt.Println("Saldos en iteracion", fromAccount.Balance, toAccount.Balance)
	}

	updateAccount1, err := store.GetAccount(context.Background(), store.db, account1.ID)
	require.NoError(t, err)

	updateAccount2, err := store.GetAccount(context.Background(), store.db, account2.ID)
	require.NoError(t, err)

	fmt.Println("Saldos finales", updateAccount1.Balance, updateAccount2.Balance)

	require.Equal(t, account1.Balance-(int64(n)*amount), updateAccount1.Balance)
	require.Equal(t, account2.Balance+(int64(n)*amount), updateAccount2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)
	
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	fmt.Println("Saldos iniciales", account1.Balance, account2.Balance)

	n := 6
	amount := int64(10)
	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 0 {
			fromAccountID = account2.ID
            toAccountID = account1.ID
        }

        go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
                FromAccountID: fromAccountID,
                ToAccountID:   toAccountID,
                Amount:        amount,
            })
			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		
		require.NoError(t, err)
	}

	updateAccount1, err := store.GetAccount(context.Background(), store.db, account1.ID)
	require.NoError(t, err)

	updateAccount2, err := store.GetAccount(context.Background(), store.db, account2.ID)
	require.NoError(t, err)

	fmt.Println("Saldos finales", updateAccount1.Balance, updateAccount2.Balance)

	require.Equal(t, account1.Balance, updateAccount1.Balance)
	require.Equal(t, account2.Balance, updateAccount2.Balance)
}
