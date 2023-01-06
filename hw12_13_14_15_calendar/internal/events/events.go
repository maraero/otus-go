package events

import "time"

type Event struct {
	ID               int64     `db:"id" json:"id"`
	Title            string    `db:"title" json:"title"`
	DateStart        time.Time `db:"date_start" json:"dateStart"`
	DateEnd          time.Time `db:"date_end" json:"dateEnd"`
	Description      string    `db:"description" json:"description"`
	UserID           string    `db:"user_id" json:"userId"`
	DateNotification time.Time `db:"date_notification" json:"dateNotification"`
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
