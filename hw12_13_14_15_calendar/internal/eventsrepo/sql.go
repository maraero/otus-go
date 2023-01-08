package eventsrepo

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/events"
)

type SQLRepository struct {
	db *sqlx.DB
}

func newSQLRepository(dbConn *sqlx.DB) *SQLRepository {
	return &SQLRepository{db: dbConn}
}

func (s *SQLRepository) CreateEvent(ctx context.Context, e events.Event) (int64, error) {
	sql := `
		INSERT INTO events (
			title,
			date_start,
			date_end,
			description,
			user_id,
			date_notification,
		)
		VALUES (
			:title,
			:date_start,
			:date_end,
			:description,
			:user_id,
			:date_notification,
		)
	`
	result, err := s.db.NamedExecContext(ctx, sql, map[string]interface{}{
		"title":             e.Title,
		"date_start":        e.DateStart,
		"date_end":          e.DateEnd,
		"description":       e.Description,
		"user_id":           e.UserID,
		"date_notification": e.DateNotification,
	})
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (s *SQLRepository) UpdateEvent(ctx context.Context, id int64, e events.Event) error {
	sql := `
		UPDATE events
		SET
			title = :title,
			date_start = :date_start,
			date_end = :date_end,
			description = :description,
		WHERE id = :id
	`
	result, err := s.db.NamedExecContext(ctx, sql, map[string]interface{}{
		"title":             e.Title,
		"date_start":        e.DateStart,
		"date_end":          e.DateEnd,
		"description":       e.Description,
		"user_id":           e.UserID,
		"date_notification": e.DateNotification,
	})
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

func (s *SQLRepository) DeleteEvent(ctx context.Context, id int64) error {
	sql := "DELETE FROM events WHERE id = :id"
	result, err := s.db.NamedExecContext(ctx, sql, map[string]interface{}{"id": id})
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

func (s *SQLRepository) GetEventListByDate(ctx context.Context, date time.Time) ([]events.Event, error) {
	year, month, day := date.Date()
	sql := `
		SELECT id, title, date_start, date_end, description, user_id, date_notification
		FROM events
		WHERE
			YEAR(date_start) = :year AND
			MONTH(date_start) = :month AND
			DAY(date_start) = :day
		ORDER BY date_start
	`
	rows, err := s.db.NamedQueryContext(ctx, sql, map[string]interface{}{
		"year":  year,
		"month": month,
		"day":   day,
	})
	if err != nil {
		return []events.Event{}, err
	}
	return parseRows(rows)
}

func (s *SQLRepository) GetEventListByWeek(ctx context.Context, date time.Time) ([]events.Event, error) {
	year, week := date.ISOWeek()
	sql := `
		SELECT id, title, date_start, date_end, description, user_id, date_notification
		FROM events
		WHERE
			YEAR(date_start) = :year AND
			WEEK(date_start) = :week
		ORDER BY date_start
	`
	rows, err := s.db.NamedQueryContext(ctx, sql, map[string]interface{}{
		"year": year,
		"week": week,
	})
	if err != nil {
		return []events.Event{}, err
	}
	return parseRows(rows)
}

func (s *SQLRepository) GetEventListByMonth(ctx context.Context, date time.Time) ([]events.Event, error) {
	year, month, _ := date.Date()
	sql := `
		SELECT id, title, date_start, date_end, description, user_id, date_notification
		FROM events
		WHERE
			YEAR(date_start) = :year AND
			MONTH(date_start) = :month
		ORDER BY date_start
	`
	rows, err := s.db.NamedQueryContext(ctx, sql, map[string]interface{}{"year": year, "month": month})
	if err != nil {
		return []events.Event{}, err
	}
	return parseRows(rows)
}

func (s *SQLRepository) GetEventByID(ctx context.Context, id int64) (events.Event, error) {
	sql := `
		SELECT id, title, date_start, date_end, description, user_id, date_notification
		FROM events
		WHERE id = $1
	`
	event := events.Event{}
	err := s.db.Get(&event, sql, id)
	if err != nil {
		return events.Event{}, err
	}
	return event, nil
}

func parseRows(rows *sqlx.Rows) ([]events.Event, error) {
	var eventList []events.Event
	for rows.Next() {
		var e events.Event
		err := rows.StructScan(e)
		if err != nil {
			return []events.Event{}, err
		}
		eventList = append(eventList, e)
	}
	return eventList, nil
}
