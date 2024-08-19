package model

import "time"

type Session struct {
	SessionID        uint32
	UserID           uint32
	HashRefreshToken string
	ExpiresAt        time.Time
}
