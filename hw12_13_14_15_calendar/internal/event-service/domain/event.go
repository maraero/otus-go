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
