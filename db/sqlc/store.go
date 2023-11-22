package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, params TransferTxParams) (TransferTxResult, error)
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db: db,
		Queries: New(db),
	}
}

func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("transaction err: %v, rollback err: %v", err, rbErr)
		}

		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID int64 `json:"to_account_id"`
	Amount int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer Transfer `json:"transfer"`
	FromAccount Account `json:"from_account"`
	ToAccount Account `json:"to_account"`
	FromEntry Entry `json:"from_entry"`
	ToEntry Entry `json:"to_entry"`
}

func (store *SQLStore) TransferTx(ctx context.Context, params TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// create transfer record
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: params.FromAccountID,
			ToAccountID: params.ToAccountID,
			Amount: params.Amount,
		})
		if err != nil {
			return err
		}

		// create from account entry
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: params.FromAccountID,
			Amount: -1 * params.Amount,
		})
		if err != nil {
			return err
		}

		// create to account entry
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: params.ToAccountID,
			Amount: params.Amount,
		})
		if err != nil {
			return err
		}

		// update accounts
		if params.FromAccountID < params.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, params.FromAccountID, -1 * params.Amount,
				params.ToAccountID, params.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, params.ToAccountID, params.Amount,
				params.FromAccountID, -1 * params.Amount)
		}

		return nil
	})

	return result, err
}

func addMoney(ctx context.Context, q *Queries, account1ID int64,
	amount1 int64, account2ID int64, amount2 int64) (account1 Account, account2 Account, err error) {
	account1, err = q.GetAccountForUpdate(ctx, account1ID)
	if err != nil {
		return
	}

	account1, err = q.UpdateAccount(ctx, UpdateAccountParams{
		ID: account1ID,
		Balance: account1.Balance + amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.GetAccountForUpdate(ctx, account2ID)
	if err != nil {
		return
	}

	account2, err = q.UpdateAccount(ctx, UpdateAccountParams{
		ID: account2ID,
		Balance: account2.Balance + amount2,
	})
	if err != nil {
		return
	}

	return
}
