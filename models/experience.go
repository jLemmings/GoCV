package models

import (
	"time"
)

type Experience struct {
	Title       string
	Description string
	From        time.Time
	To          time.Time
	Tasks       []string
}
