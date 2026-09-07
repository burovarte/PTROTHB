package domain

import "time"

type Song struct {
	ID        int64
	Name      string
	Duration  time.Duration
	CreatedAt time.Time
	UpdatedAt time.Time
}
