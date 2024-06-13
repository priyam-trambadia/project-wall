package models

import "time"

type Project struct {
	ID           int64
	Title        string
	Description  string
	OwnerID      int64
	CreatedAt    time.Time
	LastSyncedAt time.Time
}
