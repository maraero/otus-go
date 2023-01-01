package eventrepositorymemory

import (
	"sync"

	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/events"
)

type Repository struct {
	sync.RWMutex
	events map[int64]events.Event
	last   int64
}
