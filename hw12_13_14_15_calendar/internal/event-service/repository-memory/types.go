package eventrepositorymemory

import (
	"sync"

	evt "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/event-service/domain"
)

type Repository struct {
	sync.RWMutex
	events map[int64]evt.Event
	last   int64
}
