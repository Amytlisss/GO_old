package models

import "time"

type Meeting struct {
	ID        int
	UserID    int
	Date      time.Time
	Cancelled bool
	CreatedAt time.Time
	UserName  string
	UserPhone string
}
