package event

import "time"

type Event struct {
	ID               int64     `db:"id"`
	Title            string    `db:"title"`
	DateStart        time.Time `db:"date_start"`
	DateEnd          time.Time `db:"date_end"`
	Descripion       string    `db:"description"`
	UserID           string    `db:"user_id"`
	DateNotification time.Time `db:"date_notification"`
	Deleted          bool      `db:"deleted"`
}

func (e *Event) Validate() error {
	if e.Title == "" {
		return ErrEmptyTitle
	}

	if e.DateStart.IsZero() {
		return ErrEmptyDateStart
	}

	if e.DateEnd.IsZero() {
		return ErrEmptyDateEnd
	}

	if e.DateStart.After(e.DateEnd) {
		return ErrInvalidDates
	}

	if e.DateEnd.Before(time.Now()) {
		return ErrEndInThePast
	}

	return nil
}
