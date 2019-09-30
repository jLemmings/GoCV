package models

import (
	"github.com/lib/pq"
	"time"
)

type Experience struct {
	UserID      string
	ID          string `gorm:"primary_key"`
	Title       string
	Description string
	From        time.Time
	To          time.Time
	Tasks       pq.StringArray `gorm:"type:varchar(64)[]"`
}
