package dto

import "time"

type Token struct {
	ExpiresAt time.Time
	Token     string
}
