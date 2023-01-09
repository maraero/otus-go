package reporegistry

import "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/eventsrepo"

type inmemRepoRegistry struct {
	eventsRepo eventsrepo.Repository
}

func NewInMem() RepoRegistry {
	return inmemRepoRegistry{
		eventsRepo: eventsrepo.NewMemoryRepository(),
	}
}

func (r inmemRepoRegistry) GetEventsRepo() eventsrepo.Repository {
	return r.eventsRepo
}

func (r inmemRepoRegistry) DoInTransaction(txFunc InTransaction) (out any, err error) {
	return txFunc(r)
}
