package reporegistry

import "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/eventsrepo"

type InTransaction func(repoRegistry RepoRegistry) (any, error)

type RepoRegistry interface {
	GetEventsRepo() eventsrepo.Repository
	DoInTransaction(txFunc InTransaction) (out any, err error)
}
