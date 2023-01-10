package reporegistry

import (
	"database/sql"
	"fmt"

	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/eventsrepo"
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type sqlRepoRegistry struct {
	db         *sql.DB
	dbexecutor storage.DBExecutor
}

func NewSQL(db *sql.DB) RepoRegistry {
	return sqlRepoRegistry{db: db}
}

func (r sqlRepoRegistry) GetEventsRepo() eventsrepo.Repository {
	if r.dbexecutor != nil {
		return eventsrepo.NewSQLRepository(r.dbexecutor)
	}
	return eventsrepo.NewSQLRepository(r.db)
}

func (r sqlRepoRegistry) DoInTransaction(txFunc InTransaction) (out any, err error) {
	var tx *sql.Tx
	registry := r

	if r.dbexecutor == nil {
		tx, err = r.db.Begin()
		if err != nil {
			return
		}

		defer func() {
			p := recover()
			switch {
			case p != nil:
				_ = tx.Rollback()
				panic(p) // re-throw panic afterr Rollback
			case err != nil:
				xerr := tx.Rollback() // err is non-nil; don't change it
				if xerr != nil {
					err = fmt.Errorf("%w, %s", err, xerr.Error())
				}
			default:
				err = tx.Commit() // err is nil; if Commit returns errror update err
			}
		}()

		registry = sqlRepoRegistry{db: r.db, dbexecutor: tx}
	}

	out, err = txFunc(registry)
	return
}
