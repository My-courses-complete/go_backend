package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Store struct {
	*Queries
	db *pgx.Conn
}

func NewStore(db *pgx.Conn) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx error: %v, rb error: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit(ctx)
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer Transfer `json:"transfer"`
	From     Account  `json:"from"`
	To       Account  `json:"to"`
	FromEntry Entry    `json:"from_entry"`
	ToEntry   Entry    `json:"to_entry"`
}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(arg))
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.From, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:      arg.FromAccountID,
            Amount: arg.Amount,
        })
		if err!= nil {
            return err
        }

		result.To, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
            ID:      arg.ToAccountID,
            Amount: arg.Amount,
        })
		if err!= nil {
            return err
        }

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		result.To, err = q.GetAccount(ctx, arg.ToAccountID)
		if err!= nil {
            return err
        }

		return err
	})
	return result, err
}