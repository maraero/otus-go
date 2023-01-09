package reporegistry

import "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/eventsrepo"

type InTransaction func(repoRegistry RepoRegistry) (interface{}, error)

type RepoRegistry interface {
	GetPersonRepo() eventsrepo.Repository
	DoInTransaction(txFunc InTransaction) (out interface{}, err error)
}
