package database

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"
)

func RunInTx(
	ctx context.Context,
	db *bun.DB,
	opts *sql.TxOptions,
	fn func(ctx context.Context, tx bun.Tx) error,
) error {
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	done := false
	defer func() {
		if !done {
			_ = tx.Rollback()
		}
	}()

	if err := fn(ctx, tx); err != nil {
		return err
	}

	done = true
	return tx.Commit()
}
