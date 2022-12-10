package storage

import "time"

type Event struct {
	ID           string
	Title        string
	DateStart    time.Time
	DateEnd      time.Time
	Descripion   string
	UserId       string
	Notification time.Time
}
