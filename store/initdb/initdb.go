package initdb

import (
	"context"
	"database/sql"
	_ "embed"
)

//go:embed schema.sql
var ddl string

func MustInitDB(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, ddl)
	return err
}
