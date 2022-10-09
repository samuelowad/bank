package db

import (
	"context"
	"database/sql"
	"fmt"
)

//store provides all funcs to execute queries and transactions
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

//store provides all funcs to execute queries and transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}

func (s *SQLStore) execTx(ctx context.Context, fn func(queries *Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %w, rollback error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()

}

//TransferTxParams is a set of parameters for TransferTx
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

//TransferTxResult is a result of TransferTx
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entire   `json:"from_entry"`
	ToEntry     Entire   `json:"to_entry"`
}

//TransferTx performs money transfer between accounts
//it creates a transfer record and account entries and updates accounts balance with a single database transaction
func (s *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {

	var result TransferTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})

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
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		//TODO update accounts balance

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}

		return nil
	})

	if err != nil {
		return result, err
	}

	return result, nil
}

func addMoney(
	ctx context.Context,
	q *Queries,
	acc1ID int64,
	amount1 int64,
	acc2ID int64,
	amount2 int64,
) (acc1 Account, acc2 Account, err error) {
	acc1, err = q.AddAccountBalce(ctx, AddAccountBalceParams{
		ID:     acc1ID,
		Amount: amount1,
	})

	if err != nil {
		return
	}

	acc2, err = q.AddAccountBalce(ctx, AddAccountBalceParams{
		ID:     acc2ID,
		Amount: amount2,
	})

	return
}
