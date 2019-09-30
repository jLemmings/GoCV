package models

import "time"

type Education struct {
	UserID    string
	ID        string `gorm:"primary_key"`
	Title     string
	Institute string
	From      time.Time
	To        time.Time
}
