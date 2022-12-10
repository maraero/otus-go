package sqlstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	evt "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/event-service/domain"
)

type Storage struct {
	db *sqlx.DB
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context, driver string, dsn string) error {
	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	db.SetMaxOpenConns(MaxOpenConns)
	db.SetMaxIdleConns(MaxIdleConns)
	db.SetConnMaxLifetime(ConnMaxLifetime)

	s.db = db
	return s.db.PingContext(ctx)
}

func (s *Storage) Close(ctx context.Context) error {
	return s.db.Close()
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
		"description":       e.Descripion,
		"user_id":           e.UserId,
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
		"description":       e.Descripion,
		"user_id":           e.UserId,
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
	sql := "UPDATE events ыET deleted = :deleted ЦHERE id = :id"
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
	rows, err := s.db.NamedQueryContext(ctx, sql, map[string]interface{}{
		"year":  year,
		"month": month,
	})
	if err != nil {
		return []evt.Event{}, err
	}
	return parseRows(rows)
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
