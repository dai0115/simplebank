package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and trasactions
type Store struct {
	*Queries         // for query
	db       *sql.DB // for transaction
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transactions
func (s *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// generate Query object which has transaction function
	q := New(tx)
	// tx error
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err:%v rberr: %v", err, rbErr)
		}
		// in case rb was successfully completed
		return err
	}

	return tx.Commit()
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	Amount        int64 `json:"amount"`
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
}

// TransferTxResult is the result of the transfer transaction
// receive object all related to transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{}

// TransferTx performs a money transfer one account to other.
// It creates a transfer record, add account entries and update balance within a single database transaction
func (s *Store) TransferTx(ctx context.Context, args TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var err error

		txName := ctx.Value(txKey)

		fmt.Println(txName, "create transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: args.FromAccountID,
			ToAccountID:   args.ToAccountID,
			Amount:        args.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "create entry1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.FromAccountID,
			Amount:    -args.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "create entry2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.ToAccountID,
			Amount:    args.Amount,
		})
		if err != nil {
			return err
		}

		// get account -> update its balance
		fmt.Println(txName, "get account 1")
		account1, err := q.GetAccountForUpdate(ctx, args.FromAccountID)
		if err != nil {
			return err
		}

		fmt.Println(txName, "update account 1")
		result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      args.FromAccountID,
			Balance: account1.Balance - args.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "get account 2")
		account2, err := q.GetAccountForUpdate(ctx, args.ToAccountID)
		if err != nil {
			return err
		}

		fmt.Println(txName, "update account 2")
		result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      args.ToAccountID,
			Balance: account2.Balance + args.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})
	return result, err
}
