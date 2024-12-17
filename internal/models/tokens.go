package models

import "time"

type RefreshTokenRecord struct {
	ID          int
	UserID      string
	RefreshHash string
	IPAddress   string
	CreatedAt   time.Time
	ExpiresAt   time.Time
	Used        bool
}
