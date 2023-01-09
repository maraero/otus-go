package eventsrepo

import (
	"context"
	"database/sql"
	"time"

	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/events"
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/repoutils"
)

type SQLRepo struct {
	db repoutils.DBExecutor
}

func NewSQLRepository(dbConn repoutils.DBExecutor) Repository {
	return &SQLRepo{db: dbConn}
}

func (r *SQLRepo) CreateEvent(ctx context.Context, e events.Event) (int64, error) {
	sql := `
		INSERT INTO events (
			title,
			date_start,
			date_end,
			description,
			user_id,
			date_notification
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	result, err := r.db.ExecContext(ctx, sql, e.Title, e.DateStart, e.DateEnd, e.Description, e.UserID, e.DateNotification)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *SQLRepo) UpdateEvent(ctx context.Context, id int64, e events.Event) error {
	sql := `
		UPDATE events
		SET title=$1, date_start=$2, date_end=$3, description=$4
		WHERE id=$5
	`
	result, err := r.db.ExecContext(ctx, sql, e.Title, e.DateStart, e.DateEnd, e.Description, e.UserID)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return events.ErrNotFound
	}
	return nil
}

func (r *SQLRepo) DeleteEvent(ctx context.Context, id int64) error {
	sql := "DELETE FROM events WHERE id=$1"
	result, err := r.db.ExecContext(ctx, sql, id)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return events.ErrNotFound
	}
	return nil
}

func (r *SQLRepo) GetEventListByDate(ctx context.Context, date time.Time) ([]events.Event, error) {
	year, month, day := date.Date()
	sql := `
		SELECT id, title, date_start, date_end, description, user_id, date_notification
		FROM events
		WHERE
			YEAR(date_start)=$1 AND
			MONTH(date_start)=$2 AND
			DAY(date_start)=$3
		ORDER BY date_start
	`
	rows, err := r.db.QueryContext(ctx, sql, year, month, day)
	if err != nil {
		return []events.Event{}, err
	}
	defer rows.Close()
	return parseRows(rows)
}

func (r *SQLRepo) GetEventListByWeek(ctx context.Context, date time.Time) ([]events.Event, error) {
	year, week := date.ISOWeek()
	sql := `
		SELECT id, title, date_start, date_end, description, user_id, date_notification
		FROM events
		WHERE YEAR(date_start)=$1 AND WEEK(date_start)=$2
		ORDER BY date_start
	`
	rows, err := r.db.QueryContext(ctx, sql, year, week)
	if err != nil {
		return []events.Event{}, err
	}
	defer rows.Close()
	return parseRows(rows)
}

func (r *SQLRepo) GetEventListByMonth(ctx context.Context, date time.Time) ([]events.Event, error) {
	year, month, _ := date.Date()
	sql := `
		SELECT id, title, date_start, date_end, description, user_id, date_notification
		FROM events
		WHERE YEAR(date_start)=$1 AND MONTH(date_start)=$2
		ORDER BY date_start
	`
	rows, err := r.db.QueryContext(ctx, sql, year, month)
	if err != nil {
		return []events.Event{}, err
	}
	defer rows.Close()
	return parseRows(rows)
}

func (r *SQLRepo) GetEventByID(ctx context.Context, id int64) (events.Event, error) {
	sql := `
		SELECT id, title, date_start, date_end, description, user_id, date_notification
		FROM events
		WHERE id=$1
	`
	row := r.db.QueryRowContext(ctx, sql, id)
	return parseRow(row)
}

func parseRows(rows *sql.Rows) ([]events.Event, error) {
	var eventList []events.Event
	var (
		id               int64
		title            string
		dateStart        time.Time
		dateEnd          time.Time
		description      string
		userID           string
		dateNotification time.Time
	)

	for rows.Next() {
		err := rows.Scan(&id, &title, &dateStart, &dateEnd, &description, &userID, &dateNotification)
		if err != nil {
			return []events.Event{}, err
		}
		event := events.Event{
			ID:               id,
			Title:            title,
			DateStart:        dateStart,
			DateEnd:          dateEnd,
			Description:      description,
			UserID:           userID,
			DateNotification: dateNotification,
		}
		eventList = append(eventList, event)
	}

	if err := rows.Err(); err != nil {
		return []events.Event{}, err
	}

	return eventList, nil
}

func parseRow(row *sql.Row) (events.Event, error) {
	var (
		id               int64
		title            string
		dateStart        time.Time
		dateEnd          time.Time
		description      string
		userID           string
		dateNotification time.Time
	)

	err := row.Scan(&id, &title, &dateStart, &dateEnd, &description, &userID, &dateNotification)
	if err != nil {
		return events.Event{}, err
	}
	return events.Event{
		ID:               id,
		Title:            title,
		DateStart:        dateStart,
		DateEnd:          dateEnd,
		Description:      description,
		UserID:           userID,
		DateNotification: dateNotification,
	}, nil
}
