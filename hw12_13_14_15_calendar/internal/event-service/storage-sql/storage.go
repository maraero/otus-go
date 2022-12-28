package sqlstorage

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	evt "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/event-service/domain"
)

type Storage struct {
	db *sqlx.DB
}

func New(dbConn *sqlx.DB) *Storage {
	return &Storage{db: dbConn}
}

func (s *Storage) CreateEvent(ctx context.Context, e evt.Event) (int64, error) {
	sql := `
		INSERT INTO events (
			title,
			date_start,
			date_end,
			description,
			user_id,
			date_notification,
			deleted
		)
		VALUES (
			:title,
			:date_start,
			:date_end,
			:description,
			:user_id,
			:date_notification,
			:deleted
		)
	`
	result, err := s.db.NamedExecContext(ctx, sql, map[string]interface{}{
		"title":             e.Title,
		"date_start":        e.DateStart,
		"date_end":          e.DateEnd,
		"description":       e.Description,
		"user_id":           e.UserID,
		"date_notification": e.DateNotification,
		"deleted":           e.Deleted,
	})
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (s *Storage) UpdateEvent(ctx context.Context, id int64, e evt.Event) error {
	sql := `
		UPDATE events
		SET
			title = :title,
			date_start = :date_start,
			date_end = :date_end,
			description = :description,
			deleted = :deleted
		WHERE id = :id
	`
	result, err := s.db.NamedExecContext(ctx, sql, map[string]interface{}{
		"title":             e.Title,
		"date_start":        e.DateStart,
		"date_end":          e.DateEnd,
		"description":       e.Description,
		"user_id":           e.UserID,
		"date_notification": e.DateNotification,
		"deleted":           e.Deleted,
	})
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return evt.ErrNotFound
	}
	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id int64) error {
	sql := "UPDATE events SET deleted = :deleted WHERE id = :id"
	result, err := s.db.NamedExecContext(ctx, sql, map[string]interface{}{"id": id})
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return evt.ErrNotFound
	}
	return nil
}

func (s *Storage) GetEventListByDate(ctx context.Context, date time.Time) ([]evt.Event, error) {
	year, month, day := date.Date()
	sql := `
		SELECT id, title, date_start, date_end, description, user_id, date_notification
		FROM events
		WHERE
			YEAR(date_start) = :year AND
			MONTH(date_start) = :month AND
			DAY(date_start) = :day AND
			deleted = false
		ORDER BY date_start
	`
	rows, err := s.db.NamedQueryContext(ctx, sql, map[string]interface{}{
		"year":  year,
		"month": month,
		"day":   day,
	})
	if err != nil {
		return []evt.Event{}, err
	}
	return parseRows(rows)
}

func (s *Storage) GetEventListByWeek(ctx context.Context, date time.Time) ([]evt.Event, error) {
	year, week := date.ISOWeek()
	sql := `
		SELECT id, title, date_start, date_end, description, user_id, date_notification
		FROM events
		WHERE
			YEAR(date_start) = :year AND
			WEEK(date_start) = :week AND
			deleted = false
		ORDER BY date_start
	`
	rows, err := s.db.NamedQueryContext(ctx, sql, map[string]interface{}{
		"year": year,
		"week": week,
	})
	if err != nil {
		return []evt.Event{}, err
	}
	return parseRows(rows)
}

func (s *Storage) GetEventListByMonth(ctx context.Context, date time.Time) ([]evt.Event, error) {
	year, month, _ := date.Date()
	sql := `
		SELECT id, title, date_start, date_end, description, user_id, date_notification
		FROM events
		WHERE
			YEAR(date_start) = :year AND
			MONTH(date_start) = :month AND
			deleted = false
		ORDER BY date_start
	`
	rows, err := s.db.NamedQueryContext(ctx, sql, map[string]interface{}{"year": year, "month": month})
	if err != nil {
		return []evt.Event{}, err
	}
	return parseRows(rows)
}

func (s *Storage) GetEventByID(ctx context.Context, id int64) (evt.Event, error) {
	sql := `
		SELECT id, title, date_start, date_end, description, user_id, date_notification
		FROM events
		WHERE id = id
	`
	event := evt.Event{}
	err := s.db.Get(&event, sql, map[string]interface{}{"id": id})
	if err != nil {
		return evt.Event{}, err
	}
	return event, nil
}

func parseRows(rows *sqlx.Rows) ([]evt.Event, error) {
	var events []evt.Event
	for rows.Next() {
		var e evt.Event
		err := rows.StructScan(e)
		if err != nil {
			return []evt.Event{}, err
		}
		events = append(events, e)
	}
	return events, nil
}
