package reporegistry

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type DBExecutor interface {
	Get(dest interface{}, query string, args ...interface{}) error
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	NamedQueryContext(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error)
}
